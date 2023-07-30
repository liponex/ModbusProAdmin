/*
 * Copyright (C) 2023 liponex
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the terms of the GNU General Public License as published by
 * the  Free Software Foundation, either version 3 of the License, or any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.licenses/>.
 */

package serial

import (
	"log"

	"go.bug.st/serial"
)

type Serial interface {
	Opener
	Closer
	Writer
	Reader
}

type Opener interface {
	Open() (Proto, error)
}

type Closer interface {
	Close() error
}

type Writer interface {
	Write()
}

type Reader interface {
	Read()
}

type Proto struct {
	port         serial.Port
	portStr      string
	mode         *serial.Mode
	listeners    []func(proto *Proto)
	status       string
	inputBuffer  []byte
	outputBuffer []byte
}

func Open(serialPort string, mode *serial.Mode) (Proto, error) {
	port, err := serial.Open(serialPort, mode)

	proto := Proto{
		port:    port,
		portStr: serialPort,
		mode:    mode,
	}

	if err == nil {
		proto.status = "opened"
		go proto.runListeners()
	}

	return proto, err
}

func (proto *Proto) Close() error {
	err := proto.port.Close()
	if err == nil {
		proto.status = "closed"
	}
	return err
}

func (proto *Proto) Read() int {
	err := proto.port.ResetInputBuffer()
	if err != nil {
		log.Fatal(err)
	}
	n, err := proto.port.Read(proto.inputBuffer)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func (proto *Proto) AddListener(listener func(proto *Proto)) {
	proto.listeners = append(proto.listeners, listener)
}

func (proto *Proto) runListeners() {
	for proto.status != "closed" {
		if proto.inputBuffer != nil && len(proto.inputBuffer) > 0 {
			for _, listener := range proto.listeners {
				listener(proto)
			}
		}
	}
}
func (proto *Proto) String() string {
	return proto.portStr
}

func GetPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	return ports
}
