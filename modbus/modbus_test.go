package modbus

import (
	"testing"
)

func TestModbusPack(t *testing.T) {
	proto := MbPacketProto{
		0x01,
		MbPacketData{
			0x03,
			0x0001,
			0x0001,
			0x00,
			[]uint16{},
		},
		0xFFFF,
		[]uint8{},
	}
	proto.Pack()
	want := []uint8{0x01, 0x03, 0x00, 0x01, 0x00, 0x01}
	if len(want) != len(proto.PacketBuffer) {
		t.Fatal("proto.Pack() = ", proto.PacketBuffer, ", want match for ", want)
	}
	for i := 0; i < len(proto.PacketBuffer); i++ {
		if want[i] != proto.PacketBuffer[i] {
			t.Fatal("proto.Pack() = ", proto.PacketBuffer, ", want match for ", want)
		}
	}
}

func TestModbusCrc(t *testing.T) {
	proto := MbPacketProto{
		0x01,
		MbPacketData{
			0x03,
			0x0001,
			0x0001,
			0x00,
			[]uint16{},
		},
		0xFFFF,
		[]uint8{0x01, 0x03, 0x00, 0x01, 0x00, 0x01},
	}
	proto.CRC16()
	var want uint16 = 0xD5CA
	if want != proto.Crc {
		t.Fatal("proto.crc = ", proto.Crc, ", want match for ", want)
	}
}

func TestModbusReadOne(t *testing.T) {
	proto := MbPacketProto{
		0x01,
		MbPacketData{
			0x03,
			0x0001,
			0x0001,
			0x00,
			[]uint16{},
		},
		0xFFFF,
		[]uint8{},
	}
	proto.Pack()
	proto.CRC16()

	want := []uint8{0x01, 0x03, 0x00, 0x01, 0x00, 0x01, 0xD5, 0xCA}
	if len(want) != len(proto.PacketBuffer) {
		t.Fatal("proto.PacketBuffer = ", proto.PacketBuffer, ", want match for ", want)
	}
	for i := 0; i < len(proto.PacketBuffer); i++ {
		if want[i] != proto.PacketBuffer[i] {
			t.Fatal("proto.PacketBuffer = ", proto.PacketBuffer, ", want match for ", want)
		}
	}
}

func TestModbusWriteMany(t *testing.T) {
	proto := MbPacketProto{
		0x01,
		MbPacketData{
			0x10,
			0x0002,
			0x0003,
			0x06,
			[]uint16{0x000B, 0x000C, 0x000D},
		},
		0xFFFF,
		[]uint8{},
	}
	proto.Pack()
	proto.CRC16()

	want := []uint8{0x01, 0x10, 0x00, 0x02, 0x00, 0x03, 0x06, 0x00, 0x0B, 0x00, 0x0C, 0x00, 0x0D, 0xE3, 0x4D}
	if len(want) != len(proto.PacketBuffer) {
		t.Fatal("proto.PacketBuffer = ", proto.PacketBuffer, ", want match for ", want)
	}
	for i := 0; i < len(proto.PacketBuffer); i++ {
		if want[i] != proto.PacketBuffer[i] {
			t.Fatal("proto.PacketBuffer = ", proto.PacketBuffer, ", want match for ", want)
		}
	}
}
