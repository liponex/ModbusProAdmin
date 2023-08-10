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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xLayout "modbus-pro-admin/fyne/layout"
)

var (
	clientTabContainer = container.NewVBox()
	clientTabItem      = container.NewTabItem(
		"Client",
		clientTabContainer,
	)
)

func client() *container.TabItem {
	handler := container.NewVBox()

	handler.Add(NewSerialSelector())
	handler.Add(xLayout.NewResponsibleRowDistributedLayout(
		8,
		widget.NewLabel("0000"),
		widget.NewLabel("0001"),
		widget.NewLabel("0010"),
		widget.NewLabel("0011"),
		widget.NewLabel("0100"),
		widget.NewLabel("0101"),
		widget.NewLabel("0110"),
		widget.NewLabel("0111"),
		widget.NewLabel("1000"),
		widget.NewLabel("1001"),
		widget.NewLabel("1010"),
		widget.NewLabel("1011"),
		widget.NewLabel("1100"),
		widget.NewLabel("1101"),
		widget.NewLabel("1110"),
		widget.NewLabel("1111"),
	))

	clientTabContainer.Add(handler)

	return clientTabItem
}
