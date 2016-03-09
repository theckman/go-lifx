// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// MaxFrameOrigin is the max size of the Frame.Origin field.
// It only uses the top 2 bits so its maximum value is 3
const MaxFrameOrigin = ^uint8(0) >> 6

// MaxFrameProtocol is the max size of the The Frame.Protocol field.
// It only uses the top 12 bits so its maximum value is 4095.
const MaxFrameProtocol = ^uint16(0) >> 4

// FrameByteSize is the number of bytes in a marshaled Frame struct
const FrameByteSize int = 8

// ErrFrameProtocolOverflow is the error returned when the Frame.Protocol value is too large
var ErrFrameProtocolOverflow = fmt.Errorf("The Protocol field cannot be larger than %d, please choose another value (suggested: 1024)", MaxFrameProtocol)

// ErrFrameOriginOverflow is the error returned when the Frame.Origin value is too large
var ErrFrameOriginOverflow = fmt.Errorf("The Origin field cannot be larger than %d; should be set to 0", MaxFrameOrigin)

// Frame is a struct that contains some information about the message itself. This includes
// things like:
//
// 		* the size of the message
// 		* the LIFX protocol number
// 		* use of the Frame Address target field
// 		* Source identifier
type Frame struct {
	// Size of the entire message in bytes, including this field.
	Size uint16

	// Origin is the message origin indicator (must be 0)
	// Only uses the low 2 bits
	Origin uint8

	// Tagged is a boolean flag that indicates whether the FrameAddress.Target
	// field is being used to address an individual device or all devices.
	// For discovery using DeviceGetService, the Tagged field should be set to
	// true and the FrameAddress.Target should be all zeroes. The device will then
	// respond with a DeviceStateService message, which will include its own
	// MAC address in the FrameAddress.Target field. In all subsequent messages
	// that the client sends to the device, the FrameAddress.Target field should
	// be set to the device MAC address, and the tagged field should be set to false.
	//
	// The only time the Tagged field should be set to true, and the
	// FrameAddress.Target field should contain all zeros, is for the
	// DeviceGetService message.
	Tagged bool

	// Message includes a target address (must be true)
	Addressable bool

	// Protocol number; specification indicates this must be 1024
	Protocol uint16

	// Source identifier: unique value set by the client, used by responses
	Source uint32
}

// NewFrame is a function for returning a *Frame with some sane defaults.
func NewFrame() *Frame {
	return &Frame{
		Origin:      0,
		Addressable: true,
		Protocol:    1024,
	}
}

func (frame *Frame) String() string {
	if frame == nil {
		return "<*lifxprotocol.Frame(nil)>"
	}

	return fmt.Sprintf(
		"<*lifxprotocol.Frame(%p) Origin: %d, Tagged: %t, Addressable: %t, Protocol: %d, Source: 0x%x>",
		frame, frame.Origin, frame.Tagged, frame.Addressable, frame.Protocol, frame.Source,
	)
}

// MarshalPacket is a function that satisfies the Marshaler interface.
func (frame *Frame) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if frame.Origin > MaxFrameOrigin {
		return nil, ErrFrameOriginOverflow
	}

	if frame.Protocol > MaxFrameProtocol {
		return nil, ErrFrameProtocolOverflow
	}

	// TODO: enforce this in the consumer:
	// according to the protocol spec this value should be 1 (true)
	// https://github.com/LIFX/lifx-protocol-docs/blob/3de97af19703d49a97246267f08dbfd143118db0/header.md
	// if frame.Addressable == false {
	// 	frame.Addressable = true
	// }

	buf := &bytes.Buffer{}

	// write the Size field
	if err := binary.Write(buf, order, frame.Size); err != nil {
		return nil, err
	}

	// the next 16 bit value is multiple fields packed together:
	// Origin: 2
	// Tagged: 1
	// Addressable: 1
	// Protocol: 12
	mid := uint16(frame.Origin)<<14 | frame.Protocol<<4>>4

	// if Tagged set the 13th bit in mid
	if frame.Tagged {
		mid = mid | (1 << 13)
	}

	// if Addressable set the 12th bit in mid
	if frame.Addressable {
		mid = mid | (1 << 12)
	}

	// write the combination value
	if err := binary.Write(buf, order, mid); err != nil {
		return nil, err
	}

	// write the Source field
	if err := binary.Write(buf, order, frame.Source); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the Unmarshaler interface.
func (frame *Frame) UnmarshalPacket(data io.Reader, order binary.ByteOrder) error {
	if frame == nil {
		frame = &Frame{}
	}

	if err := binary.Read(data, order, &frame.Size); err != nil {
		return err
	}

	var u16 uint16

	if err := binary.Read(data, order, &u16); err != nil {
		return err
	}

	frame.Origin = uint8(u16 >> 14)    // get top 2 bits
	frame.Tagged = u16>>13&1 == 1      // get 3rd bit and eval if it's true
	frame.Addressable = u16>>12&1 == 1 // get 4th bit and eval if it's true
	frame.Protocol = u16 << 4 >> 4     // get bottom 12 bits

	if err := binary.Read(data, order, &frame.Source); err != nil {
		return err
	}

	return nil
}
