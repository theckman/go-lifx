package lifxpayloads

import (
	"bytes"
	"encoding/binary"
	"time"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestLightSetColor_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u32 uint32
	var u16 uint16
	var u8 uint8

	order := binary.LittleEndian

	lsc := &LightSetColor{
		Reserved: 20,
		Color: &LightHSBK{
			Hue:        1,
			Saturation: 2,
			Brightness: 3,
			Kelvin:     4,
		},
		Duration: 42 * time.Millisecond,
	}

	packet, err = lsc.MarshalPacket(order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Reserved
	c.Assert(binary.Read(reader, order, &u8), IsNil)
	c.Check(u8, Equals, uint8(20))

	// Color.Hue
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(1))

	// Color.Saturation
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(2))

	// Color.Brightness
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(3))

	// Color.Kelvin
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(4))

	// Duration (written as uint32 on the wire)
	c.Assert(binary.Read(reader, order, &u32), IsNil)
	c.Check(u32, Equals, uint32(42))

	//
	// Test that lsc.Duration overflow is handled gracefully
	//
	lsc.Duration = (time.Millisecond * time.Duration(^uint32(0))) + 1
	packet, err = lsc.MarshalPacket(order)
	c.Assert(err, NotNil)
	c.Check(packet, IsNil)
	c.Check(err.Error(), Equals, "LightSetColor.Duration would overflow uint32")
}

func (*TestSuite) TestLightSetColor_UnmarshalPacket(c *C) {
	var err error
	order := binary.LittleEndian
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, order, uint8(11)), IsNil)  // Reserved
	c.Assert(binary.Write(buf, order, uint16(22)), IsNil) // Color.Hue
	c.Assert(binary.Write(buf, order, uint16(33)), IsNil) // Color.Saturation
	c.Assert(binary.Write(buf, order, uint16(44)), IsNil) // Color.Brightness
	c.Assert(binary.Write(buf, order, uint16(55)), IsNil) // Color.Kelvin
	c.Assert(binary.Write(buf, order, uint32(66)), IsNil) // Duration

	lsc := &LightSetColor{}

	err = lsc.UnmarshalPacket(bytes.NewReader(buf.Bytes()), order)
	c.Assert(err, IsNil)
	c.Check(lsc.Reserved, Equals, uint8(11))
	c.Check(lsc.Color.Hue, Equals, uint16(22))
	c.Check(lsc.Color.Saturation, Equals, uint16(33))
	c.Check(lsc.Color.Brightness, Equals, uint16(44))
	c.Check(lsc.Color.Kelvin, Equals, uint16(55))
	c.Check(lsc.Duration, Equals, 66*time.Millisecond)
}

func (*TestSuite) TestLightState_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u16 uint16
	var u8 uint8

	order := binary.LittleEndian

	ls := &LightState{
		Color: &LightHSBK{
			Hue:        1,
			Saturation: 2,
			Brightness: 3,
			Kelvin:     4,
		},
		Reserved:  20,
		Power:     33,
		ReservedB: 42,
	}

	for i := 0; i < 32; i++ {
		ls.Label[i] = uint8(i + 100)
	}

	packet, err = ls.MarshalPacket(order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Color.Hue
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(1))

	// Color.Saturation
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(2))

	// Color.Brightness
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(3))

	// Color.Kelvin
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(4))

	// Reserved
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(20))

	// Power
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(33))

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Read(reader, order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+100))
	}

	// ReservedB
	c.Assert(binary.Read(reader, order, &u64), IsNil)
	c.Check(u64, Equals, uint64(42))
}

func (*TestSuite) TestLightState_UnmarshalPacket(c *C) {
	var err error
	order := binary.LittleEndian
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, order, uint16(11)), IsNil) // Color.Hue
	c.Assert(binary.Write(buf, order, uint16(22)), IsNil) // Color.Saturation
	c.Assert(binary.Write(buf, order, uint16(33)), IsNil) // Color.Brightness
	c.Assert(binary.Write(buf, order, uint16(44)), IsNil) // Color.Kelvin
	c.Assert(binary.Write(buf, order, uint16(55)), IsNil) // Reserved
	c.Assert(binary.Write(buf, order, uint16(66)), IsNil) // Power

	for i := 0; i < 32; i++ {
		c.Assert(binary.Write(buf, order, uint8(i+100)), IsNil)
	}

	c.Assert(binary.Write(buf, order, uint64(77)), IsNil) // ReservedB

	ls := &LightState{}

	err = ls.UnmarshalPacket(bytes.NewReader(buf.Bytes()), order)
	c.Assert(err, IsNil)
	c.Check(ls.Color.Hue, Equals, uint16(11))
	c.Check(ls.Color.Saturation, Equals, uint16(22))
	c.Check(ls.Color.Brightness, Equals, uint16(33))
	c.Check(ls.Color.Kelvin, Equals, uint16(44))
	c.Check(ls.Reserved, Equals, uint16(55))

	for i := 0; i < 32; i++ {
		c.Check(ls.Label[i], Equals, uint8(i+100))
	}

	c.Check(ls.ReservedB, Equals, uint64(77))
}

func (*TestSuite) TestLightSetPower_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u32 uint32
	var u16 uint16

	order := binary.LittleEndian

	lsp := &LightSetPower{
		Level:    10,
		Duration: 42 * time.Millisecond,
	}

	packet, err = lsp.MarshalPacket(order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Level
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(10))

	// Duration (written as uint32 on the wire)
	c.Assert(binary.Read(reader, order, &u32), IsNil)
	c.Check(u32, Equals, uint32(42))

	lsp.Duration = (time.Millisecond * time.Duration(^uint32(0))) + 1
	packet, err = lsp.MarshalPacket(order)
	c.Assert(err, NotNil)
	c.Check(packet, IsNil)
	c.Check(err.Error(), Equals, "LightSetPower.Duration would overflow uint32")
}

func (*TestSuite) TestLightSetPower_UnmarshalPacket(c *C) {
	var err error
	order := binary.LittleEndian
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, order, uint16(4)), IsNil)  // Level
	c.Assert(binary.Write(buf, order, uint32(22)), IsNil) // Duration

	lsp := &LightSetPower{}
	err = lsp.UnmarshalPacket(bytes.NewReader(buf.Bytes()), order)
	c.Assert(err, IsNil)
	c.Check(lsp.Level, Equals, uint16(4))
	c.Check(lsp.Duration, Equals, 22*time.Millisecond)
}

func (*TestSuite) TestLightStatePower_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u16 uint16

	order := binary.LittleEndian

	lsp := &LightStatePower{Level: 10}

	packet, err = lsp.MarshalPacket(order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Level
	c.Assert(binary.Read(reader, order, &u16), IsNil)
	c.Check(u16, Equals, uint16(10))
}

func (*TestSuite) TestLightStatePower_UnmarshalPacket(c *C) {
	var err error
	order := binary.LittleEndian
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, order, uint16(4)), IsNil) // Level

	lsp := &LightStatePower{}
	err = lsp.UnmarshalPacket(bytes.NewReader(buf.Bytes()), order)
	c.Assert(err, IsNil)
	c.Check(lsp.Level, Equals, uint16(4))
}
