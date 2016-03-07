// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"encoding/binary"
	"errors"
	"io"
)

// HeaderByteSize is the size in bytes of the header of the message.
const HeaderByteSize int = FrameByteSize + FrameAddressByteSize + ProtocolHeaderByteSize

// Header is a struct that combines the individual portions of the header
type Header struct {
	// Frame is the header component that gives some information about
	// the packet.
	Frame *Frame

	// FrameAddress is the component for defining the message target and
	// some response requirements. It's recommended to use a *FrameAddress
	// struct for this.
	FrameAddress *FrameAddress

	// ProtocolHeader is the component that contains information about
	// the payload. It's recommended to use a *ProtocolHeader struct.
	ProtocolHeader *ProtocolHeader
}

// MarshalPacket is a function that implements the Marshaler interface.
func (h *Header) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if h.Frame == nil || h.FrameAddress == nil || h.ProtocolHeader == nil {
		return nil, errors.New("none of the fields in the struct can be nil")
	}

	frame, err := h.Frame.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	frameAddress, err := h.FrameAddress.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	protocolHeader, err := h.ProtocolHeader.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	// allocate the full slice now and manually set the bytes in loops
	// later -- this is the most optimal way to do this
	packet := make([]byte, HeaderByteSize)

	fraStart := FrameByteSize
	fraEnd := fraStart + FrameAddressByteSize
	phStart := fraEnd
	phEnd := phStart + ProtocolHeaderByteSize

	// copy the Frame to packet buffer
	copy(packet, frame)

	// copy the FrameAddress to the packet buffer
	copy(packet[fraStart:fraEnd], frameAddress)

	// copy the ProtocolHeader to the packet buffer
	copy(packet[phStart:phEnd], protocolHeader)

	return packet, nil
}

// UnmarshalPacket is a function that satisfies the Unmarshaler interface.
func (h *Header) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = h.Frame.UnmarshalPacket(data, order); err != nil {
		return
	}

	if err = h.FrameAddress.UnmarshalPacket(data, order); err != nil {
		return
	}

	if err = h.ProtocolHeader.UnmarshalPacket(data, order); err != nil {
		return
	}

	return
}
