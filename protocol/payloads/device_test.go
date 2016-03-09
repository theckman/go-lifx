// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxpayloads

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestNewDeviceLabel(c *C) {
	var dl DeviceLabel
	var err error

	//
	// Test that setting a slice works
	//
	data := []byte("ohai there")
	dl, err = NewDeviceLabel(data)
	c.Assert(err, IsNil)

	label := string(dl[0:])
	c.Check(len(label), Equals, 32)
	c.Check(label, Equals, "ohai there\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	//
	// Test that setting a long slice works
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
	dl, err = NewDeviceLabel(data)
	c.Assert(err, IsNil)

	label = string(dl[0:])
	c.Check(len(label), Equals, 32)
	c.Check(label, Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")

	//
	// Test that too long of a slice fails
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456")
	dl, err = NewDeviceLabel(data)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "the slice cannot be larger than 32 bytes")
}

func (*TestSuite) TestNewDeviceLabelTrunc(c *C) {
	var dl DeviceLabel

	//
	// Test that setting a slice works
	//
	data := []byte("ohai there")
	dl = NewDeviceLabelTrunc(data)

	label := string(dl[0:])
	c.Check(len(label), Equals, 32)
	c.Check(label, Equals, "ohai there\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	//
	// Test that setting a long slice works
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
	dl = NewDeviceLabelTrunc(data)

	label = string(dl[0:])
	c.Check(len(label), Equals, 32)
	c.Check(label, Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")

	//
	// Test that too long of a slice truncates
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456")
	dl = NewDeviceLabelTrunc(data)

	label = string(dl[0:])
	c.Check(len(label), Equals, 32)
	c.Check(label, Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
}

func (*TestSuite) TestNewDeviceEchoPayloadTrunc(c *C) {
	var dep DeviceEchoPayload

	//
	// Test that setting a slice works
	//
	data := []byte("ohai there")
	dep = NewDeviceEchoPayloadTrunc(data)

	label := string(dep[0:])
	c.Check(len(label), Equals, 64)
	c.Check(label, Equals, "ohai there\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	//
	// Test that setting a long slice works
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	dep = NewDeviceEchoPayloadTrunc(data)

	label = string(dep[0:])
	c.Check(len(label), Equals, 64)
	c.Check(label, Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

	//
	// Test that too long of a slice truncates
	//
	data = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_^")
	dep = NewDeviceEchoPayloadTrunc(data)

	label = string(dep[0:])
	c.Check(len(label), Equals, 64)
	c.Check(label, Equals, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
}

func (t *TestSuite) TestDeviceStateService_String(c *C) {
	var str string

	dss := &DeviceStateService{
		Service: 3,
		Port:    56700,
	}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateService(%p): Service: 3, Port: 56700>", dss)

	str = dss.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateService_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u32 uint32
	var u8 uint8

	dss := &DeviceStateService{
		Service: 1,
		Port:    56700,
	}

	packet, err = dss.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 5)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, t.order, &u8), IsNil)
	c.Check(u8, Equals, uint8(1))
	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(56700))
}

func (t *TestSuite) TestDeviceStateService_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, uint8(42)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(8484)), IsNil)

	dss := &DeviceStateService{}

	err = dss.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dss.Service, Equals, uint8(42))
	c.Check(dss.Port, Equals, uint32(8484))
}

func (*TestSuite) TestDeviceStateHostInfo_String(c *C) {
	var str string

	dshi := &DeviceStateHostInfo{
		Signal: 0.1234,
		Tx:     1,
		Rx:     2,
	}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateHostInfo(%p): Signal: 0.1234000027179718, Tx: 1, Rx: 2>", dshi)

	str = dshi.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateHostInfo_MarshalPacket(c *C) {
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

	packet, err = dshi.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 14)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, t.order, &f32), IsNil)
	c.Check(f32, Equals, float32(42.0))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(22))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(44))

	c.Assert(binary.Read(reader, t.order, &i16), IsNil)
	c.Check(i16, Equals, int16(99))
}

func (t *TestSuite) TestDeviceStateHostInfo_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, float32(88)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(55)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(77)), IsNil)
	c.Assert(binary.Write(buf, t.order, int16(66)), IsNil)

	dshi := &DeviceStateHostInfo{}

	err = dshi.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dshi.Signal, Equals, float32(88))
	c.Check(dshi.Tx, Equals, uint32(55))
	c.Check(dshi.Rx, Equals, uint32(77))
	c.Check(dshi.Reserved, Equals, int16(66))
}

func (*TestSuite) TestDeviceStateHostFirmware_String(c *C) {
	var str string

	now := time.Now().UTC()

	dshf := &DeviceStateHostFirmware{
		Build:   uint64(now.UnixNano()),
		Version: 42,
	}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateHostFirmware(%p): Build: %s, Version: 42>", dshf, now.String())

	str = dshf.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateHostFirmware_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32

	dshf := &DeviceStateHostFirmware{
		Build:    100,
		Reserved: 30,
		Version:  200,
	}

	packet, err = dshf.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 20)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(100))

	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(30))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(200))
}

func (t *TestSuite) TestDeviceStateHostFirmware_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(84)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(99)), IsNil)

	dshf := &DeviceStateHostFirmware{}

	err = dshf.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dshf.Build, Equals, uint64(42))
	c.Check(dshf.Reserved, Equals, uint64(84))
	c.Check(dshf.Version, Equals, uint32(99))
}

func (*TestSuite) TestDeviceStateWifiInfo_String(c *C) {
	var str string

	dswi := &DeviceStateWifiInfo{
		Signal: 0.1234,
		Tx:     1,
		Rx:     2,
	}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateWifiInfo(%p): Signal: 0.1234000027179718, Tx: 1, Rx: 2>", dswi)

	str = dswi.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateWifiInfo_MarshalPacket(c *C) {
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

	packet, err = dswi.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 14)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, t.order, &f32), IsNil)
	c.Check(f32, Equals, float32(42.0))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(22))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(44))

	c.Assert(binary.Read(reader, t.order, &i16), IsNil)
	c.Check(i16, Equals, int16(99))
}

func (t *TestSuite) TestDeviceStateWifiInfo_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, float32(88)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(55)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(77)), IsNil)
	c.Assert(binary.Write(buf, t.order, int16(66)), IsNil)

	dswi := &DeviceStateWifiInfo{}

	err = dswi.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dswi.Signal, Equals, float32(88))
	c.Check(dswi.Tx, Equals, uint32(55))
	c.Check(dswi.Rx, Equals, uint32(77))
	c.Check(dswi.Reserved, Equals, int16(66))
}

func (*TestSuite) TestDeviceStateWifiFirmware_String(c *C) {
	var str string

	now := time.Now().UTC()

	dswf := &DeviceStateWifiFirmware{
		Build:   uint64(now.UnixNano()),
		Version: 42,
	}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateWifiFirmware(%p): Build: %s, Version: 42>", dswf, now.String())

	str = dswf.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateWifiFirmware_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u32 uint32

	dswf := &DeviceStateWifiFirmware{
		Build:    100,
		Reserved: 30,
		Version:  200,
	}

	packet, err = dswf.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)
	c.Assert(len(packet), Equals, 20)

	reader := bytes.NewReader(packet)

	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(100))

	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(30))

	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(200))
}

func (t *TestSuite) TestDeviceStateWifiFirmware_UnmarshalPacket(c *C) {
	var err error

	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, uint64(42)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint64(84)), IsNil)
	c.Assert(binary.Write(buf, t.order, uint32(99)), IsNil)

	dswf := &DeviceStateWifiFirmware{}

	err = dswf.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dswf.Build, Equals, uint64(42))
	c.Check(dswf.Reserved, Equals, uint64(84))
	c.Check(dswf.Version, Equals, uint32(99))
}

func (*TestSuite) TestDeviceStatePower_String(c *C) {
	var str string

	dsp := &DeviceStatePower{Level: 24}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStatePower(%p): Level: 24>", dsp)

	str = dsp.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStatePower_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u16 uint16

	dsp := &DeviceStatePower{Level: 42}

	packet, err = dsp.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Level
	c.Assert(binary.Read(reader, t.order, &u16), IsNil)
	c.Check(u16, Equals, uint16(42))
}

func (t *TestSuite) TestDeviceStatePower_UnmarshalPacket(c *C) {
	buf := &bytes.Buffer{}

	// Level
	c.Assert(binary.Write(buf, t.order, uint16(33)), IsNil)

	dsp := &DeviceStatePower{}

	c.Assert(dsp.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order), IsNil)
	c.Check(dsp.Level, Equals, uint16(33))
}

func (*TestSuite) TestDeviceStateLabel_String(c *C) {
	var str string

	label, err := NewDeviceLabel([]byte("test label"))
	c.Assert(err, IsNil)

	dsl := &DeviceStateLabel{Label: label}

	exp := fmt.Sprintf("<*lifxpayloads.DeviceStateLabel(%p): Label: \"test label\">", dsl)

	str = dsl.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateLabel_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u8 uint8

	label := make([]byte, 32)

	for i := 0; i < len(label); i++ {
		label[i] = uint8(i + 100)
	}

	dsl := &DeviceStateLabel{Label: NewDeviceLabelTrunc(label)}

	packet, err = dsl.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+100))
	}

}

func (t *TestSuite) TestDeviceStateLabel_UnmarshalPacket(c *C) {
	var err error
	buf := &bytes.Buffer{}

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+100)), IsNil)
	}

	dsl := &DeviceStateLabel{}

	err = dsl.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)

	for i := 0; i < 32; i++ {
		c.Check(dsl.Label[i], Equals, uint8(i+100))
	}
}

func (*TestSuite) TestDeviceStateVersion_String(c *C) {
	var str string

	dsv := &DeviceStateVersion{
		Vendor:  42,
		Product: 1,
		Version: 2,
	}

	exp := fmt.Sprintf(
		"<*lifxpayloads.DeviceStateVersion(%p): Vendor: 42, Product: 1, Version: 2>",
		dsv,
	)

	str = dsv.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateVersion_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u32 uint32

	dsv := &DeviceStateVersion{
		Vendor:  42,
		Product: 1,
		Version: 2,
	}

	packet, err = dsv.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Vendor
	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(42))

	// Product
	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(1))

	// Version
	c.Assert(binary.Read(reader, t.order, &u32), IsNil)
	c.Check(u32, Equals, uint32(2))
}

func (t *TestSuite) TestDeviceStateVersion_UnmarshalPacket(c *C) {
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, uint32(84)), IsNil) // Vendor
	c.Assert(binary.Write(buf, t.order, uint32(10)), IsNil) // Product
	c.Assert(binary.Write(buf, t.order, uint32(42)), IsNil) // Version

	dsv := &DeviceStateVersion{}

	c.Assert(dsv.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order), IsNil)
	c.Check(dsv.Vendor, Equals, uint32(84))
	c.Check(dsv.Product, Equals, uint32(10))
	c.Check(dsv.Version, Equals, uint32(42))
}

func (*TestSuite) TestDeviceStateInfo_String(c *C) {
	var str string

	now := time.Now().UTC()

	dsi := &DeviceStateInfo{
		Time:     uint64(now.UnixNano()),
		Uptime:   1,
		Downtime: 2,
	}

	exp := fmt.Sprintf(
		"<*lifxpayloads.DeviceStateInfo(%p): Time: %s, Uptime: 1, Downtime: 2>",
		dsi, now,
	)

	str = dsi.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateInfo_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64

	dsi := &DeviceStateInfo{
		Time:     42,
		Uptime:   1,
		Downtime: 2,
	}

	packet, err = dsi.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Time
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(42))

	// Uptime
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(1))

	// Downtime
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(2))
}

func (t *TestSuite) TestDeviceStateInfo_UnmarshalPacket(c *C) {
	buf := &bytes.Buffer{}

	c.Assert(binary.Write(buf, t.order, uint64(84)), IsNil) // Time
	c.Assert(binary.Write(buf, t.order, uint64(10)), IsNil) // Uptime
	c.Assert(binary.Write(buf, t.order, uint64(42)), IsNil) // Downtime

	dsi := &DeviceStateInfo{}

	c.Assert(dsi.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order), IsNil)
	c.Check(dsi.Time, Equals, uint64(84))
	c.Check(dsi.Uptime, Equals, uint64(10))
	c.Check(dsi.Downtime, Equals, uint64(42))
}

func (*TestSuite) TestDeviceStateLocation_String(c *C) {
	var str string

	locationStr := "location"
	label := []byte("test.bulb")

	var location [16]byte

	for i, val := range locationStr {
		location[i] = byte(val)
	}

	dsl := &DeviceStateLocation{
		Location:  location,
		Label:     NewDeviceLabelTrunc(label),
		UpdatedAt: 42,
	}

	exp := fmt.Sprintf(
		"<*lifxpayloads.DeviceStateLocation(%p): Location: \"location\", Label: \"test.bulb\", UpdatedAt: 42>",
		dsl,
	)

	str = dsl.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateLocation_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u8 uint8

	var location [16]byte

	for i := 0; i < len(location); i++ {
		location[i] = uint8(i + 200)
	}

	label := make([]byte, 32)

	for i := 0; i < len(label); i++ {
		label[i] = uint8(i + 100)
	}

	dsl := &DeviceStateLocation{
		Location:  location,
		Label:     NewDeviceLabelTrunc(label),
		UpdatedAt: 42,
	}

	packet, err = dsl.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Location
	for i := 0; i < 16; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+200))
	}

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+100))
	}

	// UpdatedAt
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(42))
}

func (t *TestSuite) TestDeviceStateLocation_UnmarshalPacket(c *C) {
	var err error
	buf := &bytes.Buffer{}

	// Location
	for i := 0; i < 16; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+200)), IsNil)
	}

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+100)), IsNil)
	}

	// UpdatedAt
	c.Assert(binary.Write(buf, t.order, uint64(84)), IsNil)

	dsl := &DeviceStateLocation{}

	err = dsl.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dsl.UpdatedAt, Equals, uint64(84))

	for i := 0; i < 16; i++ {
		c.Check(dsl.Location[i], Equals, uint8(i+200))
	}

	for i := 0; i < 32; i++ {
		c.Check(dsl.Label[i], Equals, uint8(i+100))
	}
}

func (*TestSuite) TestDeviceStateGroup_String(c *C) {
	var str string

	groupStr := "group"
	label := []byte("test.bulb")

	var group [16]byte

	for i, val := range groupStr {
		group[i] = byte(val)
	}

	dsg := &DeviceStateGroup{
		Group:     group,
		Label:     NewDeviceLabelTrunc(label),
		UpdatedAt: 42,
	}

	exp := fmt.Sprintf(
		"<*lifxpayloads.DeviceStateGroup(%p): Group: \"group\", Label: \"test.bulb\", UpdatedAt: 42>",
		dsg,
	)

	str = dsg.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceStateGroup_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u64 uint64
	var u8 uint8

	var group [16]byte

	for i := 0; i < len(group); i++ {
		group[i] = uint8(i + 200)
	}

	label := make([]byte, 32)

	for i := 0; i < len(label); i++ {
		label[i] = uint8(i + 100)
	}
	dsg := &DeviceStateGroup{
		Group:     group,
		Label:     NewDeviceLabelTrunc(label),
		UpdatedAt: 42,
	}

	packet, err = dsg.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Location
	for i := 0; i < 16; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+200))
	}

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+100))
	}

	// UpdatedAt
	c.Assert(binary.Read(reader, t.order, &u64), IsNil)
	c.Check(u64, Equals, uint64(42))
}

func (t *TestSuite) TestDeviceStateGroup_UnmarshalPacket(c *C) {
	var err error
	buf := &bytes.Buffer{}

	// Location
	for i := 0; i < 16; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+200)), IsNil)
	}

	// Label
	for i := 0; i < 32; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+100)), IsNil)
	}

	// UpdatedAt
	c.Assert(binary.Write(buf, t.order, uint64(84)), IsNil)

	dsg := &DeviceStateGroup{}

	err = dsg.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)
	c.Check(dsg.UpdatedAt, Equals, uint64(84))

	for i := 0; i < 16; i++ {
		c.Check(dsg.Group[i], Equals, uint8(i+200))
	}

	for i := 0; i < 32; i++ {
		c.Check(dsg.Label[i], Equals, uint8(i+100))
	}
}

func (*TestSuite) TestDeviceEcho_String(c *C) {
	var str string

	payload := []byte("test echo payload")

	de := &DeviceEcho{Payload: NewDeviceEchoPayloadTrunc(payload)}

	exp := fmt.Sprintf(
		"<*lifxpayloads.DeviceEcho(%p): Payload: \"test echo payload\">",
		de,
	)

	str = de.String()
	c.Check(str, Equals, exp)
}

func (t *TestSuite) TestDeviceEcho_MarshalPacket(c *C) {
	var packet []byte
	var err error
	var u8 uint8

	payload := make([]byte, 64)

	for i := 0; i < len(payload); i++ {
		payload[i] = uint8(i + 100)
	}

	de := &DeviceEcho{Payload: NewDeviceEchoPayloadTrunc(payload)}

	packet, err = de.MarshalPacket(t.order)
	c.Assert(err, IsNil)
	c.Assert(packet, NotNil)

	reader := bytes.NewReader(packet)

	// Payload
	for i := 0; i < 32; i++ {
		c.Assert(binary.Read(reader, t.order, &u8), IsNil)
		c.Check(u8, Equals, uint8(i+100))
	}

}

func (t *TestSuite) TestDeviceEcho_UnmarshalPacket(c *C) {
	var err error
	buf := &bytes.Buffer{}

	// Payload
	for i := 0; i < 64; i++ {
		c.Assert(binary.Write(buf, t.order, uint8(i+100)), IsNil)
	}

	de := &DeviceEcho{}

	err = de.UnmarshalPacket(bytes.NewReader(buf.Bytes()), t.order)
	c.Assert(err, IsNil)

	for i := 0; i < 64; i++ {
		c.Check(de.Payload[i], Equals, uint8(i+100))
	}
}
