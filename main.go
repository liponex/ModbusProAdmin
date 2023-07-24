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

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"

	"modbus-pro-admin/modbus"

	"go.bug.st/serial"
)

func main() {
	a := app.New()
	w := a.NewWindow("Modbus Pro Admin")

	settingsMenu := fyne.NewMenu("Modbus Pro Admin",
		fyne.NewMenuItem("Вырезать", func() {
			fmt.Println("Правка > Вырезать")
		}),
		fyne.NewMenuItem("Копировать", func() {
			fmt.Println("Правка > Копировать")
		}),
		fyne.NewMenuItem("Вставить", func() {
			fmt.Println("Правка > Вставить")
		}),
	)

	mainMenu := fyne.NewMainMenu(
		settingsMenu,
	)

	var serialMode = &serial.Mode{
		BaudRate: 9600,
	}

	data := modbus.MbPacketProto{
		SlaveAddr: 0x01,
		Data: modbus.MbPacketData{
			Function:   0x03,
			RegAddr:    0x0001,
			RegsAmount: 0x0001,
			DataAmount: 0x00,
		},
	}

	data.Pack()
	data.CRC16()

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	fmt.Println("Found ports:")
	for _, port := range ports {
		fmt.Printf("%v\t", port)
	}

	content := container.NewVBox(
		container.NewAppTabs(
			container.NewTabItem(
				"Tab 1",
				widget.NewSelect(
					append(ports),
					func(port string) {
						fmt.Printf("%v\t", port)
						SerialOpen(port, serialMode)
					},
				),
			),
		),
	)
	w.SetContent(content)
	w.SetMainMenu(mainMenu)
	w.SetCloseIntercept(
		func() {
			//TODO: Close serial ports on close
			w.Close()
		})

	w.ShowAndRun()
}
