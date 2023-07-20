package main

import "github.com/sigurn/crc16"

type mbPacketData struct {
	function   uint8
	regAddr    uint16
	regsAmount uint16
	dataAmount uint8
	dataBuf    []uint16
}

type mbPacketProto struct {
	slaveAddr    uint8
	data         mbPacketData
	crc          uint16
	PacketBuffer []uint8
}

func (proto *mbPacketProto) Pack() {
	proto.PacketBuffer = []uint8{
		proto.slaveAddr,
		proto.data.function,

		/* uint8(proto.data.regAddr & 0xFF),
		uint8(proto.data.regAddr >> 8),
		uint8(proto.data.regsAmount & 0xFF),
		uint8(proto.data.regsAmount >> 8), */

		uint8(proto.data.regAddr >> 8),
		uint8(proto.data.regAddr & 0xFF),
		uint8(proto.data.regsAmount >> 8),
		uint8(proto.data.regsAmount & 0xFF),
	}
	if proto.data.dataAmount > 0 {
		proto.PacketBuffer = append(proto.PacketBuffer, proto.data.dataAmount)
		for _, v := range proto.data.dataBuf {
			// proto.packetBuffer = append(proto.packetBuffer, uint8(v&0xFF), uint8(v>>8))
			proto.PacketBuffer = append(proto.PacketBuffer, uint8(v>>8), uint8(v&0xFF))
		}
	}
}

func (proto *mbPacketProto) CRC16Calculate() uint16 {
	table := crc16.MakeTable(crc16.CRC16_MODBUS)
	var crc = crc16.Checksum(proto.PacketBuffer, table)
	tmp_crc := (crc << 8) + (crc >> 8)

	proto.crc = tmp_crc
	proto.PacketBuffer = append(proto.PacketBuffer, byte(proto.crc>>8), byte(proto.crc&0xFF))
	return tmp_crc
}
