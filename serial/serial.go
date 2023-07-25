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
	Open() Proto
}

type Closer interface {
	Close()
}

type Writer interface {
	Write() Data
}

type Reader interface {
	Read() Data
}

type Proto struct {
	port serial.Port
	mode *serial.Mode
}

type Data struct {
	currentData []byte
}

func Open(serialPort string, mode *serial.Mode) (Proto, error) {
	port, err := serial.Open(serialPort, mode)

	return Proto{
		port: port,
		mode: mode,
	}, err
}

func (proto *Proto) Close() error {
	return proto.port.Close()
}

func GetPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	return ports
}
