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

func (*TestSuite) Test_phTypeToString(c *C) {
	c.Check(phTypetoString(DeviceGetService), Equals, "lifxprotocol.DeviceGetService")
	c.Check(phTypetoString(DeviceStateService), Equals, "lifxprotocol.DeviceStateService")
	c.Check(phTypetoString(DeviceGetHostInfo), Equals, "lifxprotocol.DeviceGetHostInfo")
	c.Check(phTypetoString(DeviceStateHostInfo), Equals, "lifxprotocol.DeviceStateHostInfo")
	c.Check(phTypetoString(DeviceGetHostFirmware), Equals, "lifxprotocol.DeviceGetHostFirmware")
	c.Check(phTypetoString(DeviceStateHostFirmware), Equals, "lifxprotocol.DeviceStateHostFirmware")
	c.Check(phTypetoString(DeviceGetWifiInfo), Equals, "lifxprotocol.DeviceGetWifiInfo")
	c.Check(phTypetoString(DeviceStateWifiInfo), Equals, "lifxprotocol.DeviceStateWifiInfo")
	c.Check(phTypetoString(DeviceGetWifiFirmware), Equals, "lifxprotocol.DeviceGetWifiFirmware")
	c.Check(phTypetoString(DeviceStateWifiFirmware), Equals, "lifxprotocol.DeviceStateWifiFirmware")
	c.Check(phTypetoString(DeviceGetPower), Equals, "lifxprotocol.DeviceGetPower")
	c.Check(phTypetoString(DeviceSetPower), Equals, "lifxprotocol.DeviceSetPower")
	c.Check(phTypetoString(DeviceStatePower), Equals, "lifxprotocol.DeviceStatePower")
	c.Check(phTypetoString(DeviceGetLabel), Equals, "lifxprotocol.DeviceGetLabel")
	c.Check(phTypetoString(DeviceSetLabel), Equals, "lifxprotocol.DeviceSetLabel")
	c.Check(phTypetoString(DeviceStateLabel), Equals, "lifxprotocol.DeviceStateLabel")
	c.Check(phTypetoString(DeviceGetVersion), Equals, "lifxprotocol.DeviceGetVersion")
	c.Check(phTypetoString(DeviceStateVersion), Equals, "lifxprotocol.DeviceStateVersion")
	c.Check(phTypetoString(DeviceGetInfo), Equals, "lifxprotocol.DeviceGetInfo")
	c.Check(phTypetoString(DeviceStateInfo), Equals, "lifxprotocol.DeviceStateInfo")
	c.Check(phTypetoString(DeviceAcknowledgement), Equals, "lifxprotocol.DeviceAcknowledgement")
	c.Check(phTypetoString(DeviceGetLocation), Equals, "lifxprotocol.DeviceGetLocation")
	c.Check(phTypetoString(DeviceStateLocation), Equals, "lifxprotocol.DeviceStateLocation")
	c.Check(phTypetoString(DeviceGetGroup), Equals, "lifxprotocol.DeviceGetGroup")
	c.Check(phTypetoString(DeviceStateGroup), Equals, "lifxprotocol.DeviceStateGroup")
	c.Check(phTypetoString(DeviceEchoRequest), Equals, "lifxprotocol.DeviceEchoRequest")
	c.Check(phTypetoString(DeviceEchoResponse), Equals, "lifxprotocol.DeviceEchoResponse")
	c.Check(phTypetoString(LightGet), Equals, "lifxprotocol.LightGet")
	c.Check(phTypetoString(LightSetColor), Equals, "lifxprotocol.LightSetColor")
	c.Check(phTypetoString(LightState), Equals, "lifxprotocol.LightState")
	c.Check(phTypetoString(LightGetPower), Equals, "lifxprotocol.LightGetPower")
	c.Check(phTypetoString(LightSetPower), Equals, "lifxprotocol.LightSetPower")
	c.Check(phTypetoString(LightStatePower), Equals, "lifxprotocol.LightStatePower")
	c.Check(phTypetoString(^uint16(0)), Equals, "UnknownType")
}

func (t *TestSuite) TestProtocolHeaderDeviceTypes(c *C) {
	c.Check(DeviceGetService, Equals, uint16(2))
	c.Check(DeviceStateService, Equals, uint16(3))
	c.Check(DeviceGetHostInfo, Equals, uint16(12))
	c.Check(DeviceStateHostInfo, Equals, uint16(13))
	c.Check(DeviceGetHostFirmware, Equals, uint16(14))
	c.Check(DeviceStateHostFirmware, Equals, uint16(15))
	c.Check(DeviceGetWifiInfo, Equals, uint16(16))
	c.Check(DeviceStateWifiInfo, Equals, uint16(17))
	c.Check(DeviceGetWifiFirmware, Equals, uint16(18))
	c.Check(DeviceStateWifiFirmware, Equals, uint16(19))
	c.Check(DeviceGetPower, Equals, uint16(20))
	c.Check(DeviceSetPower, Equals, uint16(21))
	c.Check(DeviceStatePower, Equals, uint16(22))
	c.Check(DeviceGetLabel, Equals, uint16(23))
	c.Check(DeviceSetLabel, Equals, uint16(24))
	c.Check(DeviceStateLabel, Equals, uint16(25))
	c.Check(DeviceGetVersion, Equals, uint16(32))
	c.Check(DeviceStateVersion, Equals, uint16(33))
	c.Check(DeviceGetInfo, Equals, uint16(34))
	c.Check(DeviceStateInfo, Equals, uint16(35))
	c.Check(DeviceAcknowledgement, Equals, uint16(45))
	c.Check(DeviceGetLocation, Equals, uint16(48))
	c.Check(DeviceStateLocation, Equals, uint16(50))
	c.Check(DeviceGetGroup, Equals, uint16(51))
	c.Check(DeviceStateGroup, Equals, uint16(53))
	c.Check(DeviceEchoRequest, Equals, uint16(58))
	c.Check(DeviceEchoResponse, Equals, uint16(59))
}

func (t *TestSuite) TestProtocolHeaderLightTypes(c *C) {
	c.Check(LightGet, Equals, uint16(101))
	c.Check(LightSetColor, Equals, uint16(102))
	c.Check(LightState, Equals, uint16(107))
	c.Check(LightGetPower, Equals, uint16(116))
	c.Check(LightSetPower, Equals, uint16(117))
	c.Check(LightStatePower, Equals, uint16(118))
}

func (*TestSuite) TestProtocolHeader_String(c *C) {
	var str string

	ph := &ProtocolHeader{Type: 2}

	exp := fmt.Sprintf(
		"<*lifxprotocol.ProtocolHeader(%p): Type: 2 (lifxprotocol.DeviceGetService)>",
		ph,
	)

	str = ph.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestProtocolHeader_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u16 uint16

	//
	// Test that Marshaling works
	//
	ph := &ProtocolHeader{
		Reserved:    1,
		Type:        2,
		ReservedEnd: 3,
	}

	packet, err = ph.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, ProtocolHeaderByteSize)

	reader := bytes.NewReader(packet)

	// read the first reserved block
	err = binary.Read(reader, t.order, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(1))

	// read the type field
	err = binary.Read(reader, t.order, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(2))

	// read the second reserved block
	err = binary.Read(reader, t.order, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(3))

	//
	// Test that Marshaling works with different inputs
	//
	ph = &ProtocolHeader{
		Reserved:    100,
		Type:        42,
		ReservedEnd: 3000,
	}

	packet, err = ph.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, ProtocolHeaderByteSize)

	reader = bytes.NewReader(packet)

	// read the first reserved block
	err = binary.Read(reader, t.order, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(100))

	// read the type field
	err = binary.Read(reader, t.order, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(42))

	// read the second reserved block
	err = binary.Read(reader, t.order, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(3000))
}

func (t *TestSuite) TestProtocolHeader_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}
	c.Assert(binary.Write(buf, t.order, uint64(3)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint16(1)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(2)), IsNil)

	//
	// Test that Unmarshaling works
	//
	ph := &ProtocolHeader{}

	err = ph.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(ph.Reserved, Equals, uint64(3))
	c.Check(ph.Type, Equals, uint16(1))
	c.Check(ph.ReservedEnd, Equals, uint16(2))

	buf.Reset()
	c.Assert(binary.Write(buf, t.order, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint16(84)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(9001)), IsNil)

	//
	// Test that Unmarshaling works with different inputs
	//
	ph = &ProtocolHeader{}

	err = ph.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(ph.Reserved, Equals, uint64(42))
	c.Check(ph.Type, Equals, uint16(84))
	c.Check(ph.ReservedEnd, Equals, uint16(9001))
}
