package lifxprotocol

import (
	"bytes"
	"encoding/binary"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestHeader_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32
	var u16 uint16
	var u8 uint8

	frame := &Frame{
		Size:        8,
		Origin:      2,
		Tagged:      true,
		Addressable: false,
		Protocol:    1024,
		Source:      42,
	}

	fraddr := &FrameAddress{
		Target:        []byte{0, 0, 0, 0, 0, 0},
		ReservedBlock: [6]uint8{0, 0, 0, 0, 0, 0},
		Reserved:      10,
		AckRequired:   false,
		ResRequired:   true,
		Sequence:      42,
	}

	ph := &ProtocolHeader{
		Reserved:    1,
		Type:        2,
		ReservedEnd: 3,
	}

	header := &Header{
		Frame:          frame,
		FrameAddress:   fraddr,
		ProtocolHeader: ph,
	}

	packet, err = header.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, Not(IsNil))

	reader := bytes.NewReader(packet)

	// Read the size field
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(8))

	// Read the middle fields that are joined together
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(uint8(u16>>14), Equals, uint8(2)) // Origin
	c.Check(u16>>13&1, Equals, uint16(1))     // Tagged
	c.Check(u16>>12&1, Equals, uint16(0))     // Addressable
	c.Check(u16<<4>>4, Equals, uint16(1024))  // Protocol

	// Read the Source field
	err = binary.Read(reader, binary.LittleEndian, &u32)
	c.Assert(err, IsNil)
	c.Check(u32, Equals, uint32(42))

	// Read the target field
	err = binary.Read(reader, binary.LittleEndian, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(0))

	// Read the 6 reserved uint8 blocks
	for i := 0; i < 6; i++ {
		err = binary.Read(reader, binary.LittleEndian, &u8)
		c.Assert(err, IsNil)
		c.Check(u8, Equals, uint8(0))
	}

	// Read the single uint8 reserved block
	err = binary.Read(reader, binary.LittleEndian, &u8)
	c.Assert(err, IsNil)
	c.Check(u8>>2, Equals, uint8(10))
	c.Check(u8>>1&1, Equals, uint8(0))
	c.Check(u8&1, Equals, uint8(1))

	// Read the sequence field
	err = binary.Read(reader, binary.LittleEndian, &u8)
	c.Assert(err, IsNil)
	c.Check(u8, Equals, uint8(42))

	// read the first reserved block
	err = binary.Read(reader, binary.LittleEndian, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(1))

	// read the type field
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(2))

	// read the second reserved block
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(3))
}

func (*TestSuite) TestHeader_UnmarshalPacket(c *C) {
	var err error
	var u64 uint64
	var u8 uint8

	origin := uint16(3)
	tagged := uint16(1)
	addressable := uint16(0)
	protocol := uint16(1024)

	u16 := origin<<14 |
		tagged<<13 |
		addressable<<12 |
		protocol<<4>>4

	buf := &bytes.Buffer{}
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(8)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, u16), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(42)), IsNil)

	u64 = uint64(65)<<55 |
		uint64(66)<<47 |
		uint64(67)<<39 |
		uint64(49)<<31 |
		uint64(50)<<23 |
		uint64(51)<<15

	c.Assert(binary.Write(buf, binary.LittleEndian, u64), IsNil)

	for i := 0; i < 6; i++ {
		c.Assert(binary.Write(buf, binary.LittleEndian, uint8(i)), IsNil)
	}

	u8 = 20<<2 | 1<<1 | 1&1
	c.Assert(binary.Write(buf, binary.LittleEndian, u8), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint8(42)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(3)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(1)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(2)), IsNil)

	frame := &Frame{}
	fra := &FrameAddress{}
	ph := &ProtocolHeader{}

	header := &Header{
		Frame:          frame,
		FrameAddress:   fra,
		ProtocolHeader: ph,
	}

	err = header.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)

	c.Check(frame.Size, Equals, uint16(8))
	c.Check(frame.Origin, Equals, uint8(origin))
	c.Check(frame.Tagged, Equals, true)
	c.Check(frame.Addressable, Equals, false)
	c.Check(frame.Protocol, Equals, uint16(1024))
	c.Check(frame.Source, Equals, uint32(42))
	c.Assert(len(fra.Target), Equals, 6)
	c.Check(fra.Target[0], Equals, byte(65))
	c.Check(fra.Target[1], Equals, byte(66))
	c.Check(fra.Target[2], Equals, byte(67))
	c.Check(fra.Target[3], Equals, byte(49))
	c.Check(fra.Target[4], Equals, byte(50))
	c.Check(fra.Target[5], Equals, byte(51))
	c.Check(fra.ReservedBlock[0], Equals, uint8(0))
	c.Check(fra.ReservedBlock[1], Equals, uint8(1))
	c.Check(fra.ReservedBlock[2], Equals, uint8(2))
	c.Check(fra.ReservedBlock[3], Equals, uint8(3))
	c.Check(fra.ReservedBlock[4], Equals, uint8(4))
	c.Check(fra.ReservedBlock[5], Equals, uint8(5))
	c.Check(fra.Reserved, Equals, uint8(20))
	c.Check(fra.AckRequired, Equals, true)
	c.Check(fra.ResRequired, Equals, true)
	c.Check(fra.Sequence, Equals, uint8(42))
	c.Check(ph.Reserved, Equals, uint64(3))
	c.Check(ph.Type, Equals, uint16(1))
	c.Check(ph.ReservedEnd, Equals, uint16(2))
}
