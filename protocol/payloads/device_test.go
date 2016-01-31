// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxpayloads

import (
	"bytes"
	"encoding/binary"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestDeviceStateService_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u32 uint32
	var u8 uint8

	dss := &DeviceStateService{
		Service: 1,
		Port:    56700,
	}

	packet, err = dss.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 5)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, binary.LittleEndian, &u8), IsNil)
	c.Check(u8, Equals, uint8(1))
	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(56700))
}

func (*TestSuite) TestDeviceStateService_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, binary.LittleEndian, uint8(42)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(8484)), IsNil)

	dss := &DeviceStateService{}

	err = dss.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(dss.Service, Equals, uint8(42))
	c.Check(dss.Port, Equals, uint32(8484))
}

func (*TestSuite) TestDeviceStateHostInfo_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var f32 float32
	var u32 uint32
	var i16 int16

	dshi := &DeviceStateHostInfo{
		Signal:   42,
		Tx:       22,
		Rx:       44,
		Reserved: 99,
	}

	packet, err = dshi.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 14)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, binary.LittleEndian, &f32), IsNil)
	c.Check(f32, Equals, float32(42.0))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(22))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(44))

	c.Assert(binary.Read(reader, binary.LittleEndian, &i16), IsNil)
	c.Check(i16, Equals, int16(99))
}

func (*TestSuite) TestDeviceStateHostInfo_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, binary.LittleEndian, float32(88)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(55)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(77)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, int16(66)), IsNil)

	dshi := &DeviceStateHostInfo{}

	err = dshi.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(dshi.Signal, Equals, float32(88))
	c.Check(dshi.Tx, Equals, uint32(55))
	c.Check(dshi.Rx, Equals, uint32(77))
	c.Check(dshi.Reserved, Equals, int16(66))
}

func (*TestSuite) TestDeviceStateHostFirmware_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32

	dshf := &DeviceStateHostFirmware{
		Build:    100,
		Reserved: 30,
		Version:  200,
	}

	packet, err = dshf.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 20)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, binary.LittleEndian, &u64), IsNil)
	c.Check(u64, Equals, uint64(100))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u64), IsNil)
	c.Check(u64, Equals, uint64(30))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(200))
}

func (*TestSuite) TestDeviceStateHostFirmware_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(84)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(99)), IsNil)

	dshf := &DeviceStateHostFirmware{}

	err = dshf.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(dshf.Build, Equals, uint64(42))
	c.Check(dshf.Reserved, Equals, uint64(84))
	c.Check(dshf.Version, Equals, uint32(99))
}

func (*TestSuite) TestDeviceStateWifiInfo_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var f32 float32
	var u32 uint32
	var i16 int16

	dswi := &DeviceStateWifiInfo{
		Signal:   42,
		Tx:       22,
		Rx:       44,
		Reserved: 99,
	}

	packet, err = dswi.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 14)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, binary.LittleEndian, &f32), IsNil)
	c.Check(f32, Equals, float32(42.0))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(22))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(44))

	c.Assert(binary.Read(reader, binary.LittleEndian, &i16), IsNil)
	c.Check(i16, Equals, int16(99))
}

func (*TestSuite) TestDeviceStateWifiInfo_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, binary.LittleEndian, float32(88)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(55)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(77)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, int16(66)), IsNil)

	dswi := &DeviceStateWifiInfo{}

	err = dswi.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(dswi.Signal, Equals, float32(88))
	c.Check(dswi.Tx, Equals, uint32(55))
	c.Check(dswi.Rx, Equals, uint32(77))
	c.Check(dswi.Reserved, Equals, int16(66))
}

func (*TestSuite) TestDeviceStateWifiFirmware_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32

	dswf := &DeviceStateWifiFirmware{
		Build:    100,
		Reserved: 30,
		Version:  200,
	}

	packet, err = dswf.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 20)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, binary.LittleEndian, &u64), IsNil)
	c.Check(u64, Equals, uint64(100))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u64), IsNil)
	c.Check(u64, Equals, uint64(30))

	c.Assert(binary.Read(reader, binary.LittleEndian, &u32), IsNil)
	c.Check(u32, Equals, uint32(200))
}

func (*TestSuite) TestDeviceStateWifiFirmware_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(84)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint32(99)), IsNil)

	dswf := &DeviceStateWifiFirmware{}

	err = dswf.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(dswf.Build, Equals, uint64(42))
	c.Check(dswf.Reserved, Equals, uint64(84))
	c.Check(dswf.Version, Equals, uint32(99))
}
