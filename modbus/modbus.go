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

package modbus

import "github.com/sigurn/crc16"

type MbPacketData struct {
	Function   uint8
	RegAddr    uint16
	RegsAmount uint16
	DataAmount uint8
	DataBuf    []uint16
}

type MbPacketProto struct {
	SlaveAddr    uint8
	Data         MbPacketData
	Crc          uint16
	PacketBuffer []uint8
}

func (proto *MbPacketProto) Pack() {
	proto.PacketBuffer = []uint8{
		proto.SlaveAddr,
		proto.Data.Function,

		/* uint8(proto.data.regAddr & 0xFF),
		uint8(proto.data.regAddr >> 8),
		uint8(proto.data.regsAmount & 0xFF),
		uint8(proto.data.regsAmount >> 8), */

		uint8(proto.Data.RegAddr >> 8),
		uint8(proto.Data.RegAddr & 0xFF),
		uint8(proto.Data.RegsAmount >> 8),
		uint8(proto.Data.RegsAmount & 0xFF),
	}
	if proto.Data.DataAmount > 0 {
		proto.PacketBuffer = append(proto.PacketBuffer, proto.Data.DataAmount)
		for _, v := range proto.Data.DataBuf {
			// proto.packetBuffer = append(proto.packetBuffer, uint8(v&0xFF), uint8(v>>8))
			proto.PacketBuffer = append(proto.PacketBuffer, uint8(v>>8), uint8(v&0xFF))
		}
	}
}

func (proto *MbPacketProto) CRC16() uint16 {
	table := crc16.MakeTable(crc16.CRC16_MODBUS)
	var crc = crc16.Checksum(proto.PacketBuffer, table)

	proto.Crc = (crc << 8) + (crc >> 8)
	proto.PacketBuffer = append(proto.PacketBuffer, byte(proto.Crc>>8), byte(proto.Crc&0xFF))
	return proto.Crc
}

func (proto *MbPacketProto) Send() {

}
