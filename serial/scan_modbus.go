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

import "modbus-pro-admin/modbus"

// ScanModbus is a split function for a Scanner that returns each Modbus command as a token.
func ScanModbus(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	proto := modbus.MbPacketProto{
		PacketBuffer: data,
		Crc:          uint16(data[len(data)-2])<<8 + uint16(data[len(data)-1]),
	}
	if proto.Validate() {
		return len(data), data, nil
	}

	return 0, nil, nil
}
