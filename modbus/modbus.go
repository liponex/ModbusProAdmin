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
	Addr   uint8
	Data   MbPacketData
	Crc    uint16
	Buffer []uint8
}

func (proto *MbPacketProto) Pack() {
	proto.Buffer = []uint8{
		proto.Addr,
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
		proto.Buffer = append(proto.Buffer, proto.Data.DataAmount)
		for _, v := range proto.Data.DataBuf {
			// proto.packetBuffer = append(proto.packetBuffer, uint8(v&0xFF), uint8(v>>8))
			proto.Buffer = append(proto.Buffer, uint8(v>>8), uint8(v&0xFF))
		}
	}
}

func (proto *MbPacketProto) CRC16() uint16 {
	table := crc16.MakeTable(crc16.CRC16_MODBUS)
	flag := proto.Crc == 0
	var crc uint16
	if flag {
		crc = crc16.Checksum(proto.Buffer, table)
	} else {
		crc = crc16.Checksum(proto.Buffer[:len(proto.Buffer)-2], table)
	}

	proto.Crc = (crc << 8) + (crc >> 8)
	if flag {
		proto.Buffer = append(proto.Buffer, byte(proto.Crc>>8), byte(proto.Crc&0xFF))
	} else {
		proto.Buffer[len(proto.Buffer)-2] = byte(proto.Crc >> 8)
		proto.Buffer[len(proto.Buffer)-1] = byte(proto.Crc & 0xFF)
	}
	return proto.Crc
}

func (proto *MbPacketProto) Send() {

}

func (proto *MbPacketProto) Validate() bool {
	proto.Crc = (uint16(proto.Buffer[len(proto.Buffer)-2]) << 8) + uint16(proto.Buffer[len(proto.Buffer)-1])
	bufCrc := proto.Crc
	trueCrc := proto.CRC16()
	return bufCrc == trueCrc
}

func (proto *MbPacketProto) Parse() {
	proto.Addr = proto.Buffer[0]
	proto.Data.Function = proto.Buffer[1]
	proto.Data.RegAddr = uint16(proto.Buffer[2])<<8 + uint16(proto.Buffer[3])
	proto.Data.RegsAmount = uint16(proto.Buffer[4])<<8 + uint16(proto.Buffer[5])
	if len(proto.Buffer) == 8 {
		proto.Crc = (uint16(proto.Buffer[6]) << 8) + uint16(proto.Buffer[7])
		return
	}
	proto.Data.DataAmount = proto.Buffer[6]
	for i := 0; i < int(proto.Data.DataAmount); i++ {
		if i+8+2 > len(proto.Buffer) {
			proto.Data.DataAmount = uint8(i - 1)
			break
		}
		proto.Data.DataBuf = append(
			proto.Data.DataBuf,
			uint16(proto.Buffer[7+i])<<8+uint16(proto.Buffer[8+i]),
		)
	}

	proto.Crc = (uint16(proto.Buffer[len(proto.Buffer)-2]) << 8) + uint16(len(proto.Buffer)-1)
}
