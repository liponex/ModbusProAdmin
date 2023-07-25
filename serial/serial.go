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

type SerialProto struct {
	port serial.Port
	mode *serial.Mode
}

func SerialOpen(serialPort string, mode *serial.Mode) SerialProto {
	port, err := serial.Open(serialPort, mode)
	if err != nil {
		log.Fatal(err)
	}

	return SerialProto{
		port: port,
		mode: mode,
	}
}

func GetSerialPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	return ports
}
