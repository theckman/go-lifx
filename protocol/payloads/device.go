// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxpayloads

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// DeviceLabel is the type corresponding to how the name of a device (the label)
// is sent over the wire. This is a 32 byte array, so helper methods exist to
// convert a slice to this type. This is NOT a payload to be send with a message.
type DeviceLabel [32]byte

// NewDeviceLabel is a function that takes a byte slice and returns a DeviceLabel.
// If the length of the byte slice is greater than 32 this will return an error.
func NewDeviceLabel(data []byte) (DeviceLabel, error) {
	if len(data) > 32 {
		return [32]byte{}, errors.New("the slice cannot be larger than 32 bytes")
	}

	dl := NewDeviceLabelTrunc(data)

	return dl, nil
}

// NewDeviceLabelTrunc is a function that takes a byte slice and returns a
// DeviceLabel. If the length of the byte slice is greater than 32 this will
// truncate the remaining bytes.
func NewDeviceLabelTrunc(data []byte) DeviceLabel {
	var dl DeviceLabel

	loops := len(data)

	if loops > 32 {
		loops = 32
	}

	for i := 0; i < loops; i++ {
		dl[i] = data[i]
	}

	return dl
}

// DeviceEchoPayload is a struct representing the payload for both the
// EchoRequest and EchoResponse message types.
type DeviceEchoPayload [64]byte

// NewDeviceEchoPayloadTrunc takes a byte slice and returns the corresponding
// DeviceEchoPayload.
func NewDeviceEchoPayloadTrunc(data []byte) DeviceEchoPayload {
	var dep DeviceEchoPayload

	loops := len(data)

	if loops > len(dep) {
		loops = len(dep)
	}

	for i := 0; i < loops; i++ {
		dep[i] = data[i]
	}

	return dep
}

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

// DeviceStatePower is the struct representing the payload for the power level
// of a device. The device sends this payload if the GetPower message is sent.
// The device expects this payload for the SetPower message.
type DeviceStatePower struct {
	Level uint16
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsp *DeviceStatePower) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dsp.Level); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsp *DeviceStatePower) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	return binary.Read(data, order, &dsp.Level)
}

// DeviceStateLabel is a struct representing the payload for setting and
// receiving the device label. The device sends this payload when responding
// to GetLabel with a StateLabel message. The client sends this payloads when
// sending a SetLabel message.
type DeviceStateLabel struct {
	Label DeviceLabel
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsl *DeviceStateLabel) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	for i := 0; i < len(dsl.Label); i++ {
		if err := binary.Write(buf, order, dsl.Label[i]); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsl *DeviceStateLabel) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	for i := 0; i < len(dsl.Label); i++ {
		if err = binary.Read(data, order, &dsl.Label[i]); err != nil {
			return
		}
	}
	return
}

// DeviceStateVersion is a struct respresenting the payload a device sends
// with the StateVersion message. It provides the hardware verson for the device.
type DeviceStateVersion struct {
	// Vendor is the Vendor ID
	Vendor uint32

	// Product is the Product ID
	Product uint32

	// Version is the hardware version
	Version uint32
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsv *DeviceStateVersion) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dsv.Vendor); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dsv.Product); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dsv.Version); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsv *DeviceStateVersion) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dsv.Vendor); err != nil {
		return
	}

	if err = binary.Read(data, order, &dsv.Product); err != nil {
		return
	}

	if err = binary.Read(data, order, &dsv.Version); err != nil {
		return
	}

	return
}

// DeviceStateInfo is the struct representation of the payload for the StateInfo
// message. This message type provides time-based information of the device.
type DeviceStateInfo struct {
	// Time is the current time in nanoseconds since the UNIX epoch
	Time uint64

	// Uptime is the time since last power on in nanoseconds
	Uptime uint64

	// Downtime is the last power off length in nanoseconds (accuracy of ~5s)
	Downtime uint64
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsi *DeviceStateInfo) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, dsi.Time); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dsi.Uptime); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, dsi.Downtime); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsi *DeviceStateInfo) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &dsi.Time); err != nil {
		return
	}

	if err = binary.Read(data, order, &dsi.Uptime); err != nil {
		return
	}

	if err = binary.Read(data, order, &dsi.Downtime); err != nil {
		return
	}

	return
}

// DeviceStateLocation location is the struct representing the device's location as
// sent by the StateLocation message.
type DeviceStateLocation struct {
	Location  [16]byte
	Label     DeviceLabel
	UpdatedAt uint64
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsl *DeviceStateLocation) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	for i := 0; i < len(dsl.Location); i++ {
		if err := binary.Write(buf, order, dsl.Location[i]); err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(dsl.Label); i++ {
		if err := binary.Write(buf, order, dsl.Label[i]); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(buf, order, dsl.UpdatedAt); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsl *DeviceStateLocation) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	for i := 0; i < len(dsl.Location); i++ {
		if err = binary.Read(data, order, &dsl.Location[i]); err != nil {
			return
		}
	}

	for i := 0; i < len(dsl.Label); i++ {
		if err = binary.Read(data, order, &dsl.Label[i]); err != nil {
			return
		}
	}

	if err = binary.Read(data, order, &dsl.UpdatedAt); err != nil {
		return
	}

	return
}

// DeviceStateGroup location is the struct representing the device's group as
// sent by the StateGroup message.
type DeviceStateGroup struct {
	Group     [16]byte
	Label     DeviceLabel
	UpdatedAt uint64
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsg *DeviceStateGroup) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	for i := 0; i < len(dsg.Group); i++ {
		if err := binary.Write(buf, order, dsg.Group[i]); err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(dsg.Label); i++ {
		if err := binary.Write(buf, order, dsg.Label[i]); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(buf, order, dsg.UpdatedAt); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (dsg *DeviceStateGroup) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	for i := 0; i < len(dsg.Group); i++ {
		if err = binary.Read(data, order, &dsg.Group[i]); err != nil {
			return
		}
	}

	for i := 0; i < len(dsg.Label); i++ {
		if err = binary.Read(data, order, &dsg.Label[i]); err != nil {
			return
		}
	}

	if err = binary.Read(data, order, &dsg.UpdatedAt); err != nil {
		return
	}

	return
}

// DeviceEcho is a struct that represents the payload for both an EchoRequest
// and an EchoResponse message.
type DeviceEcho struct {
	Payload DeviceEchoPayload
}

// MarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (de *DeviceEcho) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	for i := 0; i < len(de.Payload); i++ {
		if err := binary.Write(buf, order, de.Payload[i]); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that implements the lifxprotocol.ProtocolComponent
// interface.
func (de *DeviceEcho) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	for i := 0; i < len(de.Payload); i++ {
		if err = binary.Read(data, order, &de.Payload[i]); err != nil {
			return
		}
	}
	return
}
