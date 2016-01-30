package lifxprotocol

import (
	"bytes"
	"encoding/binary"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestProtocolHeaderDeviceTypes(c *C) {
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

func (*TestSuite) TestProtocolHeaderLightTypes(c *C) {
	c.Check(LightGet, Equals, uint16(101))
	c.Check(LightSetColor, Equals, uint16(102))
	c.Check(LightState, Equals, uint16(107))
	c.Check(LightGetPower, Equals, uint16(116))
	c.Check(LightSetPower, Equals, uint16(117))
	c.Check(LightStatePower, Equals, uint16(118))
}

func (*TestSuite) TestProtocolHeader_MarshalPacket(c *C) {
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

	packet, err = ph.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, Not(IsNil))
	c.Assert(len(packet), Equals, ProtocolHeaderByteSize)

	reader := bytes.NewReader(packet)

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

	//
	// Test that Marshaling works with different inputs
	//
	ph = &ProtocolHeader{
		Reserved:    100,
		Type:        42,
		ReservedEnd: 3000,
	}

	packet, err = ph.MarshalPacket(binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Assert(packet, Not(IsNil))
	c.Assert(len(packet), Equals, ProtocolHeaderByteSize)

	reader = bytes.NewReader(packet)

	// read the first reserved block
	err = binary.Read(reader, binary.LittleEndian, &u64)
	c.Assert(err, IsNil)
	c.Check(u64, Equals, uint64(100))

	// read the type field
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(42))

	// read the second reserved block
	err = binary.Read(reader, binary.LittleEndian, &u16)
	c.Assert(err, IsNil)
	c.Check(u16, Equals, uint16(3000))
}

func (*TestSuite) TestProtocolHeader_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(3)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(1)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(2)), IsNil)

	//
	// Test that Unmarshaling works
	//
	ph := &ProtocolHeader{}

	err = ph.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(ph.Reserved, Equals, uint64(3))
	c.Check(ph.Type, Equals, uint16(1))
	c.Check(ph.ReservedEnd, Equals, uint16(2))

	buf.Reset()
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint16(84)), IsNil)
	c.Assert(binary.Write(buf, binary.LittleEndian, uint64(9001)), IsNil)

	//
	// Test that Unmarshaling works with different inputs
	//
	ph = &ProtocolHeader{}

	err = ph.UnmarshalPacket(bytes.NewReader(buf.Bytes()), binary.LittleEndian)
	c.Assert(err, IsNil)
	c.Check(ph.Reserved, Equals, uint64(42))
	c.Check(ph.Type, Equals, uint16(84))
	c.Check(ph.ReservedEnd, Equals, uint16(9001))
}
