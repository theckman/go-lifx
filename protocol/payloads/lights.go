// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxpayloads

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

// the Duration field within the protocol specification determines how long an
// action takes to finish -- this is in milliseconds and is a uint32 on the wire
// so, we need to calculate the maximum Duration value we can support in uint32:
// 49 days, 17 hours, 2 minutes, 47.295 seconds
const lightMaxDuration = time.Millisecond * time.Duration(^uint32(0))

// ErrLightColorNotSet is the error returned when the color is not set
// on the strut trying to be marshaled.
var ErrLightColorNotSet = errors.New("a *lifxpayloads.LightHSBK must be set on the Color field before marshaling")

// LightHSBK is the struct used to represent the color and color temperature
// of a light.
//
// The color is represented as an HSB (Hue, Saturation, and Brightness) value.
type LightHSBK struct {
	// Hue is range 0 to 65535
	Hue uint16

	// Saturation is a range from 0 to 65535
	Saturation uint16

	// Brightness is a range of 0 to 65535
	Brightness uint16

	// Kevin is the color temperature of the light. The lower the warmer
	// (2500) the higher the cooler (9000).
	Kelvin uint16
}

func (hsbk *LightHSBK) String() string {
	if hsbk == nil {
		return "<*lifxpayloads.LightHSBK(nil)>"
	}

	// scale hue value to 0-359
	hue := colorRange(float64(hsbk.Hue))

	// scale saturation and brightness values to 0-100
	sat := percentageRange(float64(hsbk.Saturation))
	bri := percentageRange(float64(hsbk.Brightness))

	return fmt.Sprintf(
		"<*lifxpayloads.LightHSBK(%p): Hue: %d (%d°), Saturation: %d (%d%%), Brightness: %d (%d%%), Kelvin: %d>",
		hsbk, hsbk.Hue, hue, hsbk.Saturation, sat, hsbk.Brightness, bri, hsbk.Kelvin,
	)
}

// MarshalPacket is a function that satisfies the lifxprotocol.Marshaler
// interface.
func (hsbk *LightHSBK) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, hsbk.Hue); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, hsbk.Saturation); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, hsbk.Brightness); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, hsbk.Kelvin); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the lifxprotocol.Unmarshaler
// interface.
func (hsbk *LightHSBK) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &hsbk.Hue); err != nil {
		return
	}

	if err = binary.Read(data, order, &hsbk.Saturation); err != nil {
		return
	}

	if err = binary.Read(data, order, &hsbk.Brightness); err != nil {
		return
	}

	if err = binary.Read(data, order, &hsbk.Kelvin); err != nil {
		return
	}

	return
}

// LightSetColor is the struct representing the payload sent by a client
// to change the light state.
//
// Duration is the time it takes to transition to the new state.
type LightSetColor struct {
	Reserved uint8
	Color    *LightHSBK
	Duration time.Duration
}

func (lsc *LightSetColor) String() string {
	if lsc == nil {
		return "<*lifxpayloads.LightSetColor(nil)>"
	}

	var color string

	if lsc.Color != nil {
		color = lsc.Color.String()
	} else {
		color = "<nil>"
	}

	return fmt.Sprintf(
		"<*lifxpayloads.LightSetColor(%p): Color: %s, Duration: %s>",
		lsc, color, lsc.Duration,
	)
}

// MarshalPacket is a function that satisfies the lifxprotocol.Marshaler
// interface.
func (lsc *LightSetColor) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if lsc.Color == nil {
		return nil, ErrLightColorNotSet
	}

	// if the length of the Duration would overflow uint32
	if lsc.Duration > lightMaxDuration {
		return nil, errors.New("LightSetColor.Duration would overflow uint32")
	}

	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, lsc.Reserved); err != nil {
		return nil, err
	}

	colorPacket, err := lsc.Color.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(colorPacket); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, durToMs(lsc.Duration)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the lifxprotocol.Unmarshaler
// interface.
func (lsc *LightSetColor) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &lsc.Reserved); err != nil {
		return
	}

	if lsc.Color == nil {
		lsc.Color = &LightHSBK{}
	}

	if err = lsc.Color.UnmarshalPacket(data, order); err != nil {
		return
	}

	var u32 uint32

	if err = binary.Read(data, order, &u32); err != nil {
		return
	}

	lsc.Duration = msToDur(u32)

	return
}

// LightState is the struct representing the payload sent by the device
// to provide the current light state.
type LightState struct {
	Color    *LightHSBK
	Reserved uint16

	// Power is either 0 for off or 65535 for on
	Power uint16

	// Label is the user-identifiable name for the device.
	Label DeviceLabel

	ReservedB uint64
}

func (ls *LightState) String() string {
	if ls == nil {
		return "<*lifxpayloads.LightState(nil)>"
	}

	var color string

	if ls.Color != nil {
		color = ls.Color.String()
	} else {
		color = "<nil>"
	}

	var power string

	if ls.Power == 0 {
		power = "OFF"
	} else if ls.Power == 65535 {
		power = "ON"
	}

	label := string(bytes.Trim(ls.Label[0:], "\x00"))

	return fmt.Sprintf(
		"<*lifxpayloads.LightState(%p): Color: %s, Power: %d (%s), Label: \"%s\">",
		ls, color, ls.Power, power, label,
	)
}

// MarshalPacket is a function that satisfies the lifxprotocol.Marshaler
// interface.
func (ls *LightState) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if ls.Color == nil {
		return nil, ErrLightColorNotSet
	}

	buf := &bytes.Buffer{}

	colorPacket, err := ls.Color.MarshalPacket(order)

	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(colorPacket); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, ls.Reserved); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, ls.Power); err != nil {
		return nil, err
	}

	for i := 0; i < 32; i++ {
		if err := binary.Write(buf, order, ls.Label[i]); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(buf, order, ls.ReservedB); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the lifxprotocol.Unmarshaler
// interface.
func (ls *LightState) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if ls.Color == nil {
		ls.Color = &LightHSBK{}
	}

	if err = ls.Color.UnmarshalPacket(data, order); err != nil {
		return
	}

	if err = binary.Read(data, order, &ls.Reserved); err != nil {
		return
	}

	if err = binary.Read(data, order, &ls.Power); err != nil {
		return
	}

	for i := 0; i < 32; i++ {
		if err = binary.Read(data, order, &ls.Label[i]); err != nil {
			return
		}
	}

	if err = binary.Read(data, order, &ls.ReservedB); err != nil {
		return
	}

	return
}

// LightSetPower is a struct representing the message sent by a client to
// change the light power level.
type LightSetPower struct {
	// Level must be either 0 or 65535
	Level uint16

	// Duration is the transition time for the level change.
	Duration time.Duration
}

func (lsp *LightSetPower) String() string {
	if lsp == nil { // LumpySpacePrincess
		return "<*lifxpayloads.LightSetPower(nil)>"
	}

	var level string

	if lsp.Level == 0 {
		level = "OFF"
	} else if lsp.Level == 65535 {
		level = "ON"
	}

	return fmt.Sprintf(
		"<*lifxpayloads.LightSetPower(%p): Level: %d (%s), Duration: %s>",
		lsp, lsp.Level, level, lsp.Duration,
	)
}

// MarshalPacket is a function that satisfies the lifxprotocol.Marshaler
// interface.
func (lsp *LightSetPower) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	// if the length of the Duration would overflow uint32
	if lsp.Duration > lightMaxDuration {
		return nil, errors.New("LightSetPower.Duration would overflow uint32")
	}

	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, lsp.Level); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, durToMs(lsp.Duration)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the lifxprotocol.Unmarshaler
// interface.
func (lsp *LightSetPower) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &lsp.Level); err != nil {
		return
	}

	var u32 uint32

	if err = binary.Read(data, order, &u32); err != nil {
		return
	}

	lsp.Duration = msToDur(u32)

	return
}

// LightStatePower is the struct representing a messagent sent by a device
// to provide the current power level.
type LightStatePower struct {
	Level uint16
}

func (lsp *LightStatePower) String() string {
	if lsp == nil { // LumpySpacePrincess
		return "<*lifxpayloads.LightStatePower(nil)>"
	}

	var level string

	if lsp.Level == 0 {
		level = "OFF"
	} else if lsp.Level == 65535 {
		level = "ON"
	}

	return fmt.Sprintf(
		"<*lifxpayloads.LightStatePower(%p): Level: %d (%s)>",
		lsp, lsp.Level, level,
	)
}

// MarshalPacket is a function that satisfies the lifxprotocol.Marshaler
// interface.
func (lsp *LightStatePower) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := binary.Write(buf, order, lsp.Level); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the lifxprotocol.Unmarshaler
// interface.
func (lsp *LightStatePower) UnmarshalPacket(data io.Reader, order binary.ByteOrder) error {
	return binary.Read(data, order, &lsp.Level)
}
