// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/theckman/go-lifx/protocol/payloads"
	"github.com/theckman/go-lifx/util"

	. "gopkg.in/check.v1"
)

const MaxUint32 = ^uint32(0)

type TestSuite struct {
	order  binary.ByteOrder
	source uint32
	seed   int64
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) generateSeed(c *C) {
	tmpSeed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(tmpSeed)

	t.seed = r.Int63()
	fmt.Printf("\n!!!! github.com/theckman/go-lifx/protocol testing rand seed: %d\n\n", t.seed)

	rand.Seed(t.seed)
}

func (t *TestSuite) SetUpSuite(c *C) {
	t.generateSeed(c)

	t.order = binary.LittleEndian

	// generate a pseudorandom value between 0 and the largest uint32
	t.source = uint32(rand.Int63n(int64(MaxUint32)))
}

func (t *TestSuite) TestPacket_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32
	var u16 uint16
	var u8 uint8

	hwaddr, err := net.ParseMAC("01:23:45:67:89:ab")
	c.Assert(err, IsNil)

	rb := [6]uint8{40, 41, 42, 43, 44, 45}
	pl := [64]byte{0, 1, 2, 3, 4, 5, 6, 7}

	p := &Packet{
		Header: &Header{
			Frame: &Frame{
				Origin:      3,
				Tagged:      true,
				Addressable: true,
				Protocol:    1024,
				Source:      t.source,
			},
			FrameAddress: &FrameAddress{
				Target:        hwaddr,
				ReservedBlock: rb,
				Reserved:      50,
				AckRequired:   true,
				ResRequired:   true,
				Sequence:      128,
			},
			ProtocolHeader: &ProtocolHeader{
				Reserved:    200,
				Type:        DeviceEchoResponse,
				ReservedEnd: 2020,
			},
		},
		Payload: &lifxpayloads.DeviceEcho{
			Payload: pl,
		},
	}

	packet, err = p.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(len(packet), Equals, HeaderByteSize+64)

	reader := bytes.NewReader(packet)

	//
	// Test Header.Frame marshaling
	//
	c.Assert(binary.Read(reader, t.order, &u16), IsNil)
	c.Check(u16, Equals, uint16(HeaderByteSize+64)) // Frame.Size

	c.Assert(binary.Read(reader, t.order, &u16), IsNil)
	c.Check(uint8(u16>>14), Equals, uint8(3)) // Frame.Origin
	c.Check(u16>>13&1, Equals, uint16(1))     // Frame.Tagged
	c.Check(u16>>12&1, Equals, uint16(1))     // Frame.Addressable
	c.Check(u16<<4>>4, Equals, uint16(1024))  // Frame.Protocol

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, t.source) // Frame.Source

	//
	// Test Header.FrameAddress marshaling
	//
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	hw := lifxutil.Uint64ToHardwareAddr(u64)
	c.Check(hw.String(), Equals, hwaddr.String()) // FrameAdress.Target

	for i := range rb { // FrameAddress.ReservedBlock
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, rb[i])
	}

	c.Assert(binary.Read(reader, t.order, &u8), IsNil)
	c.Check(u8>>2, Equals, uint8(50))    // FrameAddress.Reserved
	c.Check((u8>>1)&1, Equals, uint8(1)) // FrameAddress.AckRequired
	c.Check(u8&1, Equals, uint8(1))      // FrameAddress.ResRequired

	c.Assert(binary.Read(reader, t.order, &u8), IsNil)
	c.Check(u8, Equals, uint8(128)) // FrameAddress.Sequence

	//
	// Test Header.ProtocolHeader marshaling
	//
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(200)) // ProtocolHeader.Reserved

	c.Assert(binary.Read(reader, t.order, &u16), IsNil)
	c.Check(u16, Equals, DeviceEchoResponse) // ProtocolHeader.Type

	c.Assert(binary.Read(reader, t.order, &u16), IsNil)
	c.Check(u16, Equals, uint16(2020)) // ProtocolHeader.ReservedEnd

	//
	// Test Payload marshaling
	//
	for i := range pl {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, pl[i])
	}
}

func (t *TestSuite) TestPacket_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	origin := uint16(3)
	tagged := uint16(1)
	addressable := uint16(1)
	protocol := uint16(1024)

	u16 := origin<<14 |
		tagged<<13 |
		addressable<<12 |
		protocol<<4>>4

	reserved := uint8(32)
	ackRequired := uint8(1)
	resRequired := uint8(1)

	u8 := reserved<<2 |
		(ackRequired&1)<<1 | (resRequired & 1)

	rb := [6]uint8{40, 41, 42, 43, 44, 45}

	hwaddr, err := net.ParseMAC("01:23:45:67:89:ab")
	c.Assert(err, IsNil)

	mac := lifxutil.HardwareAddrToUint64(hwaddr)

	// Packet.Header.Frame.Size
	c.Assert(binary.Write(buf, t.order, uint16(HeaderByteSize+24)), IsNil)

	// Packet.Header.Frame.
	//		Origin
	//		Tagged
	// 		Addressable
	// 		Protocol
	c.Assert(binary.Write(buf, t.order, u16), IsNil)

	// Packet.Header.Frame.Source
	c.Assert(binary.Write(buf, t.order, uint32(1122)), IsNil)

	// Packet.Header.FrameAddress.Target
	c.Assert(binary.Write(buf, t.order, mac), IsNil)

	// Packet.Header.FrameAddress.ReservedBlock
	for _, val := range rb {
		c.Assert(binary.Write(buf, t.order, val), IsNil)
	}

	// Packet.Header.FrameAddress.
	//		Reserved
	//		AckRequired
	//		ResRequired
	c.Assert(binary.Write(buf, t.order, u8), IsNil)

	// Packet.Header.FrameAdddress.Sequence
	c.Assert(binary.Write(buf, t.order, uint8(42)), IsNil)

	// Packet.Header.ProtocolHeader.Reserved
	c.Assert(binary.Write(buf, t.order, uint64(10101)), IsNil)

	// Packet.Header.ProtocolHeader.Type
	c.Assert(binary.Write(buf, t.order, DeviceStateInfo), IsNil)

	// Packet.Header.ProtocolHeader.ReservedEnd
	c.Assert(binary.Write(buf, t.order, uint16(2424)), IsNil)

	// Packet.Payload (lifxpayloads.DeviceStateInfo)
	c.Assert(binary.Write(buf, t.order, uint64(11223344)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(22334455)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(33445566)), IsNil)

	reader := bytes.NewReader(buf.Bytes())

	p := &Packet{}

	err = p.UnmarshalPacket(reader, t.order)
	c.Assert(err, IsNil)

	c.Check(p.Header.Frame.Size, Equals, uint16(HeaderByteSize+24))
	c.Check(p.Header.Frame.Origin, Equals, uint8(origin))
	c.Check(p.Header.Frame.Tagged, Equals, true)
	c.Check(p.Header.Frame.Addressable, Equals, true)
	c.Check(p.Header.Frame.Protocol, Equals, protocol)
	c.Check(p.Header.FrameAddress.Target.String(), Equals, hwaddr.String())
	c.Check(p.Header.FrameAddress.ReservedBlock, DeepEquals, rb)
	c.Check(p.Header.FrameAddress.Reserved, Equals, uint8(32))
	c.Check(p.Header.FrameAddress.AckRequired, Equals, true)
	c.Check(p.Header.FrameAddress.ResRequired, Equals, true)
	c.Check(p.Header.FrameAddress.Sequence, Equals, uint8(42))
	c.Check(p.Header.ProtocolHeader.Reserved, Equals, uint64(10101))
	c.Check(p.Header.ProtocolHeader.Type, Equals, DeviceStateInfo)
	c.Check(p.Header.ProtocolHeader.ReservedEnd, Equals, uint16(2424))

	c.Assert(p.Payload, NotNil)

	payload, ok := p.Payload.(*lifxpayloads.DeviceStateInfo)
	c.Assert(ok, Equals, true)
	c.Assert(payload, NotNil)
	c.Check(payload.Time, Equals, uint64(11223344))
	c.Check(payload.Uptime, Equals, uint64(22334455))
	c.Check(payload.Downtime, Equals, uint64(33445566))
}
