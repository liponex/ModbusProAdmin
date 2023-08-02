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

package gui

import (
	"log"
	"modbus-pro-admin/modbus"
	"modbus-pro-admin/serial"

	serialLib "go.bug.st/serial"

	xWidget "modbus-pro-admin/fyne/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	OpenSerials []serial.Proto
)

func serialSelector() *fyne.Container {
	var (
		serialPorts = append(
			[]string{"Disconnect"},
			serial.GetPorts()...,
		)
		serialMode = &serialLib.Mode{
			BaudRate: 9600,
			DataBits: 8,
		}
	)

	labelSerialPort := widget.NewLabel("Serial port:")
	selectSerialPort := xWidget.NewSelectWithPrevLink[serial.Proto](
		serialPorts,
		func(hasPrev *bool, prevSelected *serial.Proto, newString string) {
			if len(serialPorts) == 0 {
				return
			}
			if newString == prevSelected.String() {
				return
			}
			if newString == "Disconnect" {
				if !*hasPrev {
					return
				}
				if prevSelected != nil {
					prevSelected.Terminate()
				}
				*hasPrev = false
				return
			}

			proto := new(serial.Proto)
			proto, err := serial.Open(newString, serialMode)
			if err != nil {
				log.Fatal("Can't open serial new ", newString)
			}

			proto.AddListener(func(proto *serial.Proto) {
				flag := true
				word := []uint8{1, 3, 0, 2, 0, 2, 101, 203}
				i := 0
				for i < len(proto.InputBuffer) && i < len(word) {
					if proto.InputBuffer[i] != word[i] {
						flag = false
						break
					}
					i++
				}
				if flag {
					packet := modbus.MbPacketProto{
						PacketBuffer: []uint8{1, 3, 4, 0, 250, 0, 150},
					}
					packet.CRC16()
					proto.OutputBuffer = packet.PacketBuffer
					err := proto.Write()
					if err != nil {
						log.Fatalln(err)
					}
				}
				//fmt.Println(proto.InputBuffer)
			})

			if !*hasPrev {
				*prevSelected = *proto
				*hasPrev = true
				return
			}
			prevSelected.Terminate()
			*prevSelected = *proto
		},
	)
	selectSerialPort.PlaceHolder = "(Select port)"
	buttonUpdate := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		serialPorts = serial.GetPorts()
		if len(serialPorts) > 0 {
			selectSerialPort.Options = append(
				[]string{"Disconnect"},
				serialPorts...,
			)
		} else {
			selectSerialPort.Options = []string{"No ports available"}
		}
	})

	return container.NewHBox(
		labelSerialPort,
		selectSerialPort,
		buttonUpdate,
	)
}
