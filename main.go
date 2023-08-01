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
	"modbus-pro-admin/gui"
	"modbus-pro-admin/modbus"

	"fyne.io/fyne/v2/data/binding"
)

var (
	PortsBinding = binding.BindStringList(
		&[]string{},
	)
)

func main() {
	data := modbus.MbPacketProto{
		Addr: 0x01,
		Data: modbus.MbPacketData{
			Function:   0x03,
			RegAddr:    0x0000,
			RegsAmount: 0x0001,
			DataAmount: 0x01,
			DataBuf: []uint16{
				0x000F,
			},
		},
	}

	data.Pack()
	data.CRC16()

	gui.Gui()
}
