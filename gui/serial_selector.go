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
	"errors"
	"fmt"
	"log"
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
		func(hasPrev *bool, prevSelected *serial.Proto, new string) {
			if len(serialPorts) == 0 {
				return
			}
			if new == prevSelected.String() {
				return
			}
			if new == "Disconnect" {
				if !*hasPrev {
					return
				}
				go func(port serial.Proto) {
					var err = errors.New("")
					for err != nil {
						err = port.Close()
					}
					fmt.Println("Port", port.String(), "closed")
				}(*prevSelected)
				*hasPrev = false
				return
			}

			fmt.Println(new)
			proto, err := serial.Open(new, serialMode)
			if err != nil {
				log.Fatal("Can't open serial new ", new)
			}

			proto.AddListener(func(proto *serial.Proto) {
				fmt.Println(proto.InputBuffer)
			})
			if !*hasPrev {
				*prevSelected = proto
				*hasPrev = true
				return
			}
			go func(port serial.Proto) {
				var err = errors.New("")
				for err != nil {
					err = port.Close()
				}
				fmt.Println("Port", port.String(), "closed")
			}(*prevSelected)
			*prevSelected = proto
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
