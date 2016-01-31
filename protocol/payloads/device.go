// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxpayloads

import (
	"bytes"
	"encoding/binary"
	"io"
)

// DeviceStateService is the response to the DeviceGetService message.
//
// Provides the device Service and port. If the Service is temporarily
// unavailable, then the port value will be 0.
type DeviceStateService struct {
	// Service describes the type of service exposed by the device.
	// 		1: UDP
	Service uint8

	// Port is the port the device is listening on the network. For
	// compatibility reasons it's recommended that clients bind to port
	// 56700.
	Port uint32
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dss *DeviceStateService) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dss.Service); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dss.Port); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dss *DeviceStateService) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dss.Service); err != nil {
		return
	}

	if err = binary.Read(data, order, &dss.Port); err != nil {
		return
	}

	return
}

// DeviceStateHostInfo is the response to the DeviceGetHostInfo message.
// It provides host MCU information.
type DeviceStateHostInfo struct {
	// Signal is the radio receive signal strength in milliwatts.
	Signal float32

	// Tx is the number of bytes transmitted since power on.
	Tx uint32

	// Rx is the number of bytes received since power on.
	Rx uint32

	Reserved int16
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dshi *DeviceStateHostInfo) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dshi.Signal); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dshi.Tx); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dshi.Rx); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dshi.Reserved); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dshi *DeviceStateHostInfo) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dshi.Signal); err != nil {
		return
	}

	if err = binary.Read(data, order, &dshi.Tx); err != nil {
		return
	}

	if err = binary.Read(data, order, &dshi.Rx); err != nil {
		return
	}

	if err = binary.Read(data, order, &dshi.Reserved); err != nil {
		return
	}

	return
}

// DeviceStateHostFirmware is the response to the DeviceGetHosFirmware message.
// This provides information about the host's firmware.
type DeviceStateHostFirmware struct {
	// Build is the firmware build time (absolute time in nanoseconds since epoch).
	Build uint64

	Reserved uint64

	// Version is the firmware version of the host.
	Version uint32
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dshf *DeviceStateHostFirmware) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dshf.Build); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dshf.Reserved); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dshf.Version); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dshf *DeviceStateHostFirmware) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dshf.Build); err != nil {
		return
	}

	if err = binary.Read(data, order, &dshf.Reserved); err != nil {
		return
	}

	if err = binary.Read(data, order, &dshf.Version); err != nil {
		return
	}

	return
}

// DeviceStateWifiInfo is the response to the DeviceGetWifiInfo message.
// It provides Wifi subsystem information.
type DeviceStateWifiInfo struct {
	// Signal is the radio receive signal strength in milliwatts
	Signal float32

	// Tx is the number of bytes transmitted since power on.
	Tx uint32

	// Rx is the nimber of bytes received since power on.
	Rx uint32

	Reserved int16
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dswi *DeviceStateWifiInfo) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dswi.Signal); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dswi.Tx); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dswi.Rx); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dswi.Reserved); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dswi *DeviceStateWifiInfo) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dswi.Signal); err != nil {
		return
	}

	if err = binary.Read(data, order, &dswi.Tx); err != nil {
		return
	}

	if err = binary.Read(data, order, &dswi.Rx); err != nil {
		return
	}

	if err = binary.Read(data, order, &dswi.Reserved); err != nil {
		return
	}

	return
}

// DeviceStateWifiFirmware is the response to the GetWifiFirmware message.
// This provides Wifi subsystem information.
type DeviceStateWifiFirmware struct {
	// Build is the firmware build time (absolute time in nanoseconds since epoch)
	Build uint64

	Reserved uint64

	// Version is the subsystem firmware version.
	Version uint32
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dswf *DeviceStateWifiFirmware) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dswf.Build); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dswf.Reserved); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dswf.Version); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dswf *DeviceStateWifiFirmware) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dswf.Build); err != nil {
		return
	}

	if err = binary.Read(data, order, &dswf.Reserved); err != nil {
		return
	}

	if err = binary.Read(data, order, &dswf.Version); err != nil {
		return
	}

	return
}
