// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// Package lifxprotocol is used for marshaling and unmarshaling different LIFX
// protocol packets to and from the wire, respectively. This package is not
// meant to be consumed by those wanting to interface with their LIFX devices
// in Golang. This package is designed to be used by the LIFX Golang library
// for communicating with devices. Users are meant to consume that package
// (which unfortunately doesn't exist, yet).
//
// This package uses the lifpayloads sub-package to generate payloads for
// the individual packets.
package lifxprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/theckman/go-lifx/protocol/payloads"
)

const maxUint16 = int(^uint16(0))

// Marshaler is the interface for Marshaling packets.
// The order parameter can either be binary.LittleEndian or binary.BigEndian.
// The LIFX protocol uses little-endian encoding at the time of writing.
type Marshaler interface {
	MarshalPacket(order binary.ByteOrder) ([]byte, error)
}

// Unmarshaler is the interface for Unmarshaling packets.
// The data parameter is an io.Reader that contains the contents of the packet.
// The order parameter can either be binary.LittleEndian or binary.BigEndian.
// The LIFX protocol uses little-endian encoding at the time of writing.
type Unmarshaler interface {
	UnmarshalPacket(data io.Reader, order binary.ByteOrder) error
}

// PacketComponent is the interface used for each component to implement.
type PacketComponent interface {
	Marshaler
	Unmarshaler
	fmt.Stringer
}

// Packet is the struct for an individual message whether inbound or
// outbound.
type Packet struct {
	// Header is the component of a packet that details who it is for,
	// what it contains, and how the recipient is meant to reply to
	// the packet.
	Header *Header

	// Payload is the message payload. The contents of the payload vary
	// based on what is being sent. The recipient knows how to parse this
	// message based on the value of the Header.ProtocolHeader.Type field.
	Payload PacketComponent
}

func (p *Packet) String() string {
	if p == nil {
		return "<*lifxprotocol.Packet(nil)>"
	}

	var hStr, pStr string

	if p.Header == nil {
		hStr = "<nil>"
	} else {
		hStr = p.Header.String()
	}

	if p.Payload == nil {
		pStr = "<nil>"
	} else {
		pStr = p.Payload.String()
	}

	return fmt.Sprintf(
		"<*lifxprotocol.Packet(%p): Header: %s, Payload: %s>",
		p, hStr, pStr,
	)
}

// MarshalPacket is a function that satisfies the Marshaler interface.
func (p *Packet) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if p.Header == nil {
		return nil, errors.New("the Header field cannot be nil")
	}

	if p.Payload == nil {
		return nil, errors.New("the Payload field cannot be nil")
	}

	payload, err := p.Payload.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	// calculate the total byte size
	tbs := HeaderByteSize + len(payload)

	// check for overflow of Header.Frame.Size field
	if tbs > maxUint16 {
		return nil, fmt.Errorf("size of packet (%d) would overflow Packet.Header.Frame.Size uint16 field (max %d)", tbs, maxUint16)
	}

	// we now know how big the message is now, so let's set it
	p.Header.Frame.Size = uint16(tbs)

	header, err := p.Header.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	packet := make([]byte, tbs)

	// copy the header to the beginning of the packet
	copy(packet, header)

	// copy the payload immediately following the header
	copy(packet[HeaderByteSize:tbs], payload)

	return packet, nil
}

func packetComponentByType(t uint16) PacketComponent {
	switch t {
	case DeviceStateService:
		return &lifxpayloads.DeviceStateService{}

	case DeviceStateHostInfo:
		return &lifxpayloads.DeviceStateHostInfo{}

	case DeviceStateHostFirmware:
		return &lifxpayloads.DeviceStateHostFirmware{}

	case DeviceStateWifiInfo:
		return &lifxpayloads.DeviceStateWifiInfo{}

	case DeviceStateWifiFirmware:
		return &lifxpayloads.DeviceStateWifiFirmware{}

	case DeviceStatePower, DeviceSetPower:
		return &lifxpayloads.DeviceStatePower{}

	case DeviceStateLabel, DeviceSetLabel:
		return &lifxpayloads.DeviceStateLabel{}

	case DeviceStateVersion:
		return &lifxpayloads.DeviceStateVersion{}

	case DeviceStateInfo:
		return &lifxpayloads.DeviceStateInfo{}

	case DeviceStateLocation:
		return &lifxpayloads.DeviceStateInfo{}

	case DeviceStateGroup:
		return &lifxpayloads.DeviceStateGroup{}

	case DeviceEchoResponse, DeviceEchoRequest:
		return &lifxpayloads.DeviceEcho{}

	case LightSetColor:
		return &lifxpayloads.LightSetColor{}

	case LightState:
		return &lifxpayloads.LightState{}

	case LightSetPower:
		return &lifxpayloads.LightSetPower{}

	case LightStatePower:
		return &lifxpayloads.LightStatePower{}

	default:
		return nil
	}
}

func (p *Packet) unmarshalPayload(data io.Reader, order binary.ByteOrder) (PacketComponent, error) {
	if p.Header.ProtocolHeader == nil {
		return nil, errors.New("the ProtocolHeader cannot be nil")
	}

	var pc PacketComponent

	// figure out the payload type so we can unmarshal it
	if pc = packetComponentByType(p.Header.ProtocolHeader.Type); pc == nil {
		return nil, errors.New("unknown message type")
	}

	if err := pc.UnmarshalPacket(data, order); err != nil {
		return nil, err
	}

	return pc, nil
}

// UnmarshalPacket is a function that implements the Unmarshaler interface.
func (p *Packet) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	hdr := &Header{
		Frame:          &Frame{},
		FrameAddress:   &FrameAddress{},
		ProtocolHeader: &ProtocolHeader{},
	}

	if err = hdr.UnmarshalPacket(data, order); err != nil {
		return
	}

	p.Header = hdr

	payload, err := p.unmarshalPayload(data, order)

	if err != nil {
		return
	}

	p.Payload = payload

	return
}
