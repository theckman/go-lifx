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
	"io"
)

// ProtocolComponent is the interface used for each component to implement.
type ProtocolComponent interface {
	MarshalPacket(order binary.ByteOrder) ([]byte, error)
	UnmarshalPacket(data io.Reader, order binary.ByteOrder) error
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
	Payload ProtocolComponent
}
