// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"

	. "gopkg.in/check.v1"
)

func (t *TestSuite) Test_NewFrameAddress(c *C) {
	var fa *FrameAddress
	fa = NewFrameAddress()
	c.Assert(fa, NotNil)
}

func (*TestSuite) TestFrameAddress_String(c *C) {
	var str string

	fraddr := &FrameAddress{
		Target:        []byte{1, 2, 3, 4, 5, 6},
		ReservedBlock: [6]uint8{0, 0, 0, 0, 0, 0},
		Reserved:      10,
		AckRequired:   false,
		ResRequired:   true,
		Sequence:      42,
	}

	exp := fmt.Sprintf(
		"<*lifxprotocol.FrameAddress(%p): Target: 01:02:03:04:05:06, AckRequired: false, ResRequired: true, Sequence: 42>",
		fraddr,
	)

	str = fraddr.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestFrameAddress_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u8 uint8

	//
	// Test that Marshaling works
	//
	fraddr := &FrameAddress{
		Target:        []byte{0, 0, 0, 0, 0, 0},
		ReservedBlock: [6]uint8{0, 0, 0, 0, 0, 0},
		Reserved:      10,
		AckRequired:   false,
		ResRequired:   true,
		Sequence:      42,
	}

	packet, err = fraddr.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Check(len(packet), Equals, FrameAddressByteSize)

	reader := bytes.NewReader(packet)

	// Read the target field
	err = binary.Read(reader, t.order, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(0))

	// Read the 6 reserved uint8 blocks
	for i := 0; i < 6; i++ {
		err = binary.Read(reader, t.order, &u8)
		c.Assert(err, IsNil)
		c.Check(u8, Equals, uint8(0))
	}

	// Read the single uint8 reserved block
	err = binary.Read(reader, t.order, &u8)
	c.Assert(err, IsNil)
	c.Check(u8>>2, Equals, uint8(10))
	c.Check(u8>>1&1, Equals, uint8(0))
	c.Check(u8&1, Equals, uint8(1))

	// Read the sequence field
	err = binary.Read(reader, t.order, &u8)
	c.Assert(err, IsNil)
	c.Check(u8, Equals, uint8(42))

	//
	// Test that Marshaling works with different fields
	//
	fraddr = &FrameAddress{
		Target:        []byte{65, 66, 67, 49, 50, 51},
		ReservedBlock: [6]uint8{1, 2, 3, 4, 5, 6},
		Reserved:      11,
		AckRequired:   true,
		ResRequired:   false,
		Sequence:      22,
	}

	packet, err = fraddr.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Check(len(packet), Equals, FrameAddressByteSize)

	reader = bytes.NewReader(packet)

	// Read the target field
	err = binary.Read(reader, t.order, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(2351197419751440384))

	// Read the 6 reserved uint8 blocks
	for i := 0; i < 6; i++ {
		err = binary.Read(reader, t.order, &u8)
		c.Assert(err, IsNil)
		c.Check(u8, Equals, uint8(i+1))
	}

	// Read the single uint8 reserved block
	err = binary.Read(reader, t.order, &u8)
	c.Assert(err, IsNil)
	c.Check(u8>>2, Equals, uint8(11))
	c.Check(u8>>1&1, Equals, uint8(1))
	c.Check(u8&1, Equals, uint8(0))

	// Read the sequence field
	err = binary.Read(reader, t.order, &u8)
	c.Assert(err, IsNil)
	c.Check(u8, Equals, uint8(22))
}

func (t *TestSuite) TestFrameAddress_UnmarshalPacket(c *C) {
	var err error
	var u64 uint64
	var u8 uint8

	buf := new(bytes.Buffer)

	u64 = uint64(65)<<55 |
		uint64(66)<<47 |
		uint64(67)<<39 |
		uint64(49)<<31 |
		uint64(50)<<23 |
		uint64(51)<<15

	err = binary.Write(buf, t.order, u64)
	c.Assert(err, IsNil)

	for i := 0; i < 6; i++ {
		err = binary.Write(buf, t.order, uint8(i))
		c.Assert(err, IsNil)
	}

	u8 = 20<<2 | 1<<1 | 1&1
	err = binary.Write(buf, t.order, u8)
	c.Assert(err, IsNil)

	err = binary.Write(buf, t.order, uint8(42))
	c.Assert(err, IsNil)

	//
	// Test that Unmarshaling works
	//
	fra := &FrameAddress{}

	err = fra.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
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

	buf = new(bytes.Buffer)

	u64 = uint64(65)<<55 |
		uint64(66)<<47 |
		uint64(67)<<39 |
		uint64(49)<<31 |
		uint64(50)<<23 |
		uint64(51)<<15

	err = binary.Write(buf, t.order, u64)
	c.Assert(err, IsNil)

	for i := 0; i < 6; i++ {
		err = binary.Write(buf, t.order, uint8(i+1))
		c.Assert(err, IsNil)
	}

	u8 = 20<<2 | 0<<1 | 0&1
	err = binary.Write(buf, t.order, u8)
	c.Assert(err, IsNil)

	err = binary.Write(buf, t.order, uint8(42))
	c.Assert(err, IsNil)

	//
	// Test that Unmarshaling works with different inputs
	//
	fra = &FrameAddress{}

	err = fra.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Assert(len(fra.Target), Equals, 6)
	c.Check(fra.Target[0], Equals, byte(65))
	c.Check(fra.Target[1], Equals, byte(66))
	c.Check(fra.Target[2], Equals, byte(67))
	c.Check(fra.Target[3], Equals, byte(49))
	c.Check(fra.Target[4], Equals, byte(50))
	c.Check(fra.Target[5], Equals, byte(51))
	c.Check(fra.ReservedBlock[0], Equals, uint8(1))
	c.Check(fra.ReservedBlock[1], Equals, uint8(2))
	c.Check(fra.ReservedBlock[2], Equals, uint8(3))
	c.Check(fra.ReservedBlock[3], Equals, uint8(4))
	c.Check(fra.ReservedBlock[4], Equals, uint8(5))
	c.Check(fra.ReservedBlock[5], Equals, uint8(6))
	c.Check(fra.Reserved, Equals, uint8(20))
	c.Check(fra.AckRequired, Equals, false)
	c.Check(fra.ResRequired, Equals, false)
	c.Check(fra.Sequence, Equals, uint8(42))
}
