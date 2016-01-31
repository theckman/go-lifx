// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"encoding/binary"
	"errors"
	"io"
)

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

// MarshalPackage is a funtion that implements the ProtocolComponent interface.
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
	packet := make([]byte, len(frame)+len(frameAddress)+len(protocolHeader))

	// pointer for the packet byte slice
	// we use it to know which byte we are setting
	var j int

	// loop over the frame bytes and put them in the packet slice
	// loop over the frameAddress bytes and put them in the packet slice
	// using the j pointer to put the bytes in the right location within
	// packet
	for i := 0; i < len(frame); i++ {
		packet[j] = frame[i]
		j++ // move pointer
	}

	// loop over the frameAddress bytes and put them in the packet slice
	for i := 0; i < len(frameAddress); i++ {
		packet[j] = frameAddress[i]
		j++
	}

	// loop over the protocolHeader bytes and put them in the packet slice
	for i := 0; i < len(protocolHeader); i++ {
		packet[j] = protocolHeader[i]
		j++
	}

	return packet, nil
}

// UnmarshalPackage is a funtion that implements the ProtocolComponent interface.
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
