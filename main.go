package main

import (
	"fmt"
)

func main() {
	fmt.Println("SuperModbus")
	data := mbPacketProto{
		0x01,
		mbPacketData{
			0x03,
			0x0001,
			0x0001,
			0x00,
			[]uint16{},
		},
		0xFFFF,
		[]uint8{},
	}

	data.Pack()
	data.CRC16Calculate()
}
