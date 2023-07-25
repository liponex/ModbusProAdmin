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
	"fmt"
	"modbus-pro-admin/serial"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func serialSelector() *fyne.Container {
	var serialPorts = serial.GetSerialPorts()

	labelSerialPort := widget.NewLabel("Serial port:")
	selectSerialPort := widget.NewSelect(
		serialPorts,
		func(port string) {
			if len(serialPorts) > 0 {
				fmt.Printf("%v\t", port)
				// SerialOpen(port, serialMode)
			} else {
			}
		},
	)
	selectSerialPort.PlaceHolder = "(Select port)"
	buttonUpdate := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		serialPorts = serial.GetSerialPorts()
		if len(serialPorts) > 0 {
			selectSerialPort.Options = serialPorts
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
