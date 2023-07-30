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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"log"
)

func Gui() {
	a := app.New()
	w := a.NewWindow("Modbus Pro Admin")

	settingsMenu := fyne.NewMenu("Modbus Pro Admin",
		fyne.NewMenuItem("Cut", func() {
			fmt.Println("Edit > Cut")
		}),
		fyne.NewMenuItem("Copy", func() {
			fmt.Println("Edit > Copy")
		}),
		fyne.NewMenuItem("Paste", func() {
			fmt.Println("Edit > Paste")
		}),
	)

	mainMenu := fyne.NewMainMenu(
		settingsMenu,
	)

	content := container.NewVBox(
		container.NewAppTabs(
			client(),
		),
	)
	w.SetContent(content)
	w.SetMainMenu(mainMenu)
	w.SetCloseIntercept(
		func() {
			for _, port := range OpenSerials {
				err := port.Close()
				if err != nil {
					log.Fatal(err)
				}
			}
			w.Close()
		})

	w.ShowAndRun()
}
