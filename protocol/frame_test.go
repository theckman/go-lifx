// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"bytes"
	"encoding/binary"

	. "gopkg.in/check.v1"
)

func (t *TestSuite) Test_NewFrame(c *C) {
	var f *Frame
	f = NewFrame()
	c.Assert(f, NotNil)
	c.Check(f.Origin, Equals, uint8(0))
	c.Check(f.Addressable, Equals, true)
	c.Check(f.Protocol, Equals, uint16(1024))
}

func (t *TestSuite) TestFrame_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u16 uint16
	var u32 uint32

	//
	// Test that Marshaling works
	//
	frame := &Frame{
		Size:        8,
		Origin:      2,
		Tagged:      true,
		Addressable: false,
		Protocol:    1024,
		Source:      42,
	}

	packet, err = frame.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Check(len(packet), Equals, FrameByteSize)

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

	//
	// Test that Marshaling works with some adjusted fields
	//
	frame = &Frame{
		Size:        10,
		Origin:      0,
		Tagged:      false,
		Addressable: true,
		Protocol:    4095,
		Source:      4242,
	}

	packet, err = frame.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Check(len(packet), Equals, FrameByteSize)

	reader = bytes.NewReader(packet)

	// Read the size field
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(10))

	// Read the middle fields that are joined together
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(uint8(u16>>14), Equals, uint8(0)) // Origin
	c.Check(u16>>13&1, Equals, uint16(0))     // Tagged
	c.Check(u16>>12&1, Equals, uint16(1))     // Addressable
	c.Check(u16<<4>>4, Equals, uint16(4095))  // Protocol

	// Read the Source field
	err = binary.Read(reader, binary.LittleEndian, &u32)
	c.Assert(err, IsNil)
	c.Check(u32, Equals, uint32(4242))

	//
	// Test that overflowing the Origin field throws an error
	//
	frame = &Frame{
		Size:        8,
		Origin:      4,
		Tagged:      false,
		Addressable: false,
		Protocol:    1024,
		Source:      42,
	}

	packet, err = frame.MarshalPacket(binary.LittleEndian)
	c.Check(err, Equals, ErrFrameOriginOverflow)
	c.Check(packet, IsNil)

	//
	// Test that overflowing the Protocol field throws an error
	//
	frame = &Frame{
		Size:        8,
		Origin:      0,
		Tagged:      false,
		Addressable: false,
		Protocol:    4096,
		Source:      42,
	}

	packet, err = frame.MarshalPacket(binary.LittleEndian)
	c.Check(err, Equals, ErrFrameProtocolOverflow)
	c.Check(packet, IsNil)
}

func (t *TestSuite) TestFrame_UnmarshalPacket(c *C) {
	var err error

	origin := uint16(3)
	tagged := uint16(1)
	addressable := uint16(0)
	protocol := uint16(1024)

	u16 := origin<<14 |
		tagged<<13 |
		addressable<<12 |
		protocol<<4>>4

	buf := new(bytes.Buffer)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(8)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, u16), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(42)), IsNil)

	//
	// Test that Unmarshaling works
	//
	frame := &Frame{}

	err = frame.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)

	c.Check(frame.Size, Equals, uint16(8))
	c.Check(frame.Origin, Equals, uint8(origin))
	c.Check(frame.Tagged, Equals, true)
	c.Check(frame.Addressable, Equals, false)
	c.Check(frame.Protocol, Equals, uint16(1024))
	c.Check(frame.Source, Equals, uint32(42))

	origin = uint16(0)
	tagged = uint16(0)
	addressable = uint16(1)
	protocol = uint16(4095)

	u16 = origin<<14 |
		tagged<<13 |
		addressable<<12 |
		protocol<<4>>4

	buf = new(bytes.Buffer)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(10)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, u16), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(4242)), IsNil)

	//
	// Test that Unmarshaling works with some adjusted fields
	//
	frame = &Frame{}

	err = frame.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)

	c.Check(frame.Size, Equals, uint16(10))
	c.Check(frame.Origin, Equals, uint8(origin))
	c.Check(frame.Tagged, Equals, false)
	c.Check(frame.Addressable, Equals, true)
	c.Check(frame.Protocol, Equals, uint16(4095))
	c.Check(frame.Source, Equals, uint32(4242))
}
