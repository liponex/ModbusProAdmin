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

	clientTabContainer.Add(xLayout.NewResponsibleRowDistributedLayout(
		3,
		newSerialSelector(),
		widget.NewLabel("123"),
		widget.NewLabel("456"),
	))

	return clientTabItem
}
