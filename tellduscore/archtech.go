package tellduscore

import (
	"bytes"
	"io"
)

// GetRawCommand encodes a tellstick command for Archtech
func GetRawCommand(house uint, unit int, method int, level int) string {
	var command bytes.Buffer

	command.WriteString("T")
	command.WriteByte(byte(127))
	command.WriteByte(byte(255))
	command.WriteByte(byte(24))
	command.WriteByte(byte(1))

	unit--

	if method == TellstickDim {
		command.WriteByte(147)
	} else {
		command.WriteByte(132)
	}

	// House
	var m bytes.Buffer
	var i int
	var ui uint
	for i = 25; i >= 0; i-- {
		ui = uint(i)
		if house&(1<<ui) == 0 {
			m.WriteString("01")
		} else {
			m.WriteString("10")
		}
	}
	m.WriteString("01")

	// on/off/dim
	if method == TellstickDim {
		m.WriteString("00")
	} else if method == TellstickTurnoff {
		m.WriteString("01")
	} else if method == TellstickTurnon {
		m.WriteString("10")
	}

	// code (unit)
	for i = 4; i >= 0; i-- {
		ui = uint(i)
		if unit&(1<<ui) == 0 {
			m.WriteString("01")
		} else {
			m.WriteString("10")
		}
	}

	// dim level
	if method == TellstickDim {
		newLevel := level / 16
		for i = 4; i >= 0; i-- {
			ui = uint(i)
			if newLevel&(1<<ui) == 0 {
				m.WriteString("01")
			} else {
				m.WriteString("10")
			}
		}
	}

	code := 9
	i = 0
	for v, err := m.ReadByte(); err != io.EOF; v, err = m.ReadByte() {
		code <<= 4
		if string(v) == "1" {
			code |= 8
		} else {
			code |= 10
		}

		if i%2 == 0 {
			command.WriteByte(byte(code))
			code = 0
		}
		i++
	}

	command.WriteString("+")

	return command.String()
}
