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
	"fyne.io/fyne/v2"
	"log"
	xLayout "modbus-pro-admin/fyne/layout"
	"modbus-pro-admin/modbus"
	"modbus-pro-admin/serial"
	"strconv"

	serialLib "go.bug.st/serial"

	xWidget "modbus-pro-admin/fyne/widget"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	OpenSerials []serial.Proto
)

type serialSelector struct {
	serialPort *serial.Proto
	serialMode *serialLib.Mode
	container  *fyne.Container
}

func (s *serialSelector) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return s.container.Layout.MinSize(objects)
}

func (s *serialSelector) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	s.container.Layout.Layout(objects, containerSize)
}

func NewSerialSelector() fyne.CanvasObject {
	varSerialSelector := new(serialSelector)
	var (
		serialPorts = append(
			[]string{"Disconnect"},
			serial.GetPorts()...,
		)
	)

	varSerialSelector.serialMode = new(serialLib.Mode)
	*varSerialSelector.serialMode = serialLib.Mode{
		BaudRate: 9600,
		Parity:   serialLib.NoParity,
		DataBits: 8,
		StopBits: serialLib.OneStopBit,
	}
	varSerialSelector.serialPort = new(serial.Proto)

	labelSerialPort := widget.NewLabel("Serial port:")
	selectSerialPort := xWidget.NewSelectWithPrevLink[serial.Proto](
		serialPorts,
		func(hasPrev *bool, prevSelected *serial.Proto, newString string) {
			if len(serialPorts) == 0 {
				return
			}
			if newString == prevSelected.String() {
				return
			}
			if newString == "Disconnect" {
				if !*hasPrev {
					return
				}
				if prevSelected != nil {
					prevSelected.Terminate()
				}
				*hasPrev = false
				return
			}

			proto := new(serial.Proto)
			proto, err := serial.Open(newString, varSerialSelector.serialMode)
			if err != nil {
				log.Fatal("Can't open serial new ", newString)
			}

			proto.AddListener(func(proto *serial.Proto) {
				flag := true
				word := []uint8{1, 3, 0, 2, 0, 2, 101, 203}
				i := 0
				for i < len(proto.InputBuffer) && i < len(word) {
					if proto.InputBuffer[i] != word[i] {
						flag = false
						break
					}
					i++
				}
				if flag {
					packet := modbus.MbPacketProto{
						Buffer: []uint8{1, 3, 4, 0, 250, 0, 150},
					}
					packet.CRC16()
					proto.OutputBuffer = packet.Buffer
					err := proto.Write()
					if err != nil {
						log.Fatalln(err)
					}
				}
			})

			if !*hasPrev {
				*prevSelected = *proto
				*varSerialSelector.serialPort = *proto
				*hasPrev = true
				return
			}
			prevSelected.Terminate()
			*prevSelected = *proto
			*varSerialSelector.serialPort = *proto
		},
	)
	selectSerialPort.PlaceHolder = "(Select port)"
	buttonUpdate := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		serialPorts = serial.GetPorts()
		if len(serialPorts) > 0 {
			selectSerialPort.Options = append(
				[]string{"Disconnect"},
				serialPorts...,
			)
		} else {
			selectSerialPort.Options = []string{"No ports available"}
		}
	})

	labelSerialBaudRate := widget.NewLabel("BaudRate: ")
	selectSerialBaudRate := widget.NewSelectEntry(serial.DefaultRates)
	selectSerialBaudRate.SetText("9600")
	selectSerialBaudRate.OnSubmitted = func(baudStr string) {
		var err error
		varSerialSelector.serialMode.BaudRate, err = strconv.Atoi(baudStr)
		if err != nil {
			varSerialSelector.serialMode.BaudRate, _ = strconv.Atoi(serial.DefaultRates[0])
		}
		selectSerialBaudRate.SetText(strconv.Itoa(varSerialSelector.serialMode.BaudRate))

		varSerialSelector.serialPort.Terminate()
		varSerialSelector.serialPort, err = serial.Open(
			varSerialSelector.serialPort.String(),
			varSerialSelector.serialMode,
		)
	}
	objects := []fyne.CanvasObject{
		container.NewHBox(
			labelSerialPort,
			selectSerialPort,
			buttonUpdate,
		),
		container.NewHBox(
			labelSerialBaudRate,
			selectSerialBaudRate,
		),
	}

	varSerialSelector.container = xLayout.NewResponsibleRowDistributedLayout(len(objects), objects...)
	return container.New(varSerialSelector, objects...)
}
