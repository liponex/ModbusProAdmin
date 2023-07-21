package main

import (
	"fmt"
	"log"

	"modbus-pro-admin/modbus"

	"go.bug.st/serial"
)

func main() {
	fmt.Println("SuperModbus")
	data := modbus.MbPacketProto{
		SlaveAddr: 0x01,
		Data: modbus.MbPacketData{
			Function:   0x03,
			RegAddr:    0x0001,
			RegsAmount: 0x0001,
			DataAmount: 0x00,
		},
	}

	data.Pack()
	data.CRC16()

	serialMode := &serial.Mode{
		BaudRate: 9600,
	}
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	fmt.Println("Found ports:")
	for _, port := range ports {
		fmt.Printf("%v\t", port)
	}
	SerialOpen(ports[2], serialMode)
}
