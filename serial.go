package main

import (
	"log"

	"go.bug.st/serial"
)

type SerialProto struct {
	port serial.Port
	mode *serial.Mode
}

func SerialOpen(serialPort string, mode *serial.Mode) SerialProto {
	port, err := serial.Open(serialPort, mode)
	if err != nil {
		log.Fatal(err)
	}

	return SerialProto{
		port: port,
		mode: mode,
	}
}
