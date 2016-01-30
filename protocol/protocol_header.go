package lifxprotocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ProtocolHeaderByteSize is the number of bytes in a marshaled packet.
const ProtocolHeaderByteSize int = 12

// These values are for use in the Type field. They define the type of
// message within the payload of the packet. This group of values are for
// generic device messages.
const (
	DeviceGetService        uint16 = 2
	DeviceStateService      uint16 = 3
	DeviceGetHostInfo       uint16 = 12
	DeviceStateHostInfo     uint16 = 13
	DeviceGetHostFirmware   uint16 = 14
	DeviceStateHostFirmware uint16 = 15
	DeviceGetWifiInfo       uint16 = 16
	DeviceStateWifiInfo     uint16 = 17
	DeviceGetWifiFirmware   uint16 = 18
	DeviceStateWifiFirmware uint16 = 19
	DeviceGetPower          uint16 = 20
	DeviceSetPower          uint16 = 21
	DeviceStatePower        uint16 = 22
	DeviceGetLabel          uint16 = 23
	DeviceSetLabel          uint16 = 24
	DeviceStateLabel        uint16 = 25
	DeviceGetVersion        uint16 = 32
	DeviceStateVersion      uint16 = 33
	DeviceGetInfo           uint16 = 34
	DeviceStateInfo         uint16 = 35
	DeviceAcknowledgement   uint16 = 45
	DeviceGetLocation       uint16 = 48
	DeviceStateLocation     uint16 = 50
	DeviceGetGroup          uint16 = 51
	DeviceStateGroup        uint16 = 53
	DeviceEchoRequest       uint16 = 58
	DeviceEchoResponse      uint16 = 59
)

// These values are for use in the Type field. They define the type of
// message within the payload of the packet. This group of values are for
// device messages specific to LIFX lightbulbs.
const (
	LightGet        uint16 = 101
	LightSetColor   uint16 = 102
	LightState      uint16 = 107
	LightGetPower   uint16 = 116
	LightSetPower   uint16 = 117
	LightStatePower uint16 = 118
)

// ProtocolHeader is a struct that contains information about the payload contents
// (i.e., what actions to take)
type ProtocolHeader struct {
	// Reserved is reserved according to the protocol documentation
	Reserved uint64

	// Type is the message type used to determine the message payload
	Type uint16

	// ReservedEnd is additional reserved space as defined by the protocol
	// documentation
	ReservedEnd uint16
}

func (ph *ProtocolHeader) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	buf := &bytes.Buffer{}

	// write the first reserved block
	if err := binary.Write(buf, order, ph.Reserved); err != nil {
		return nil, err
	}

	// write the type field, which indicates payload type
	if err := binary.Write(buf, order, ph.Type); err != nil {
		return nil, err
	}

	// write the last reserved block
	if err := binary.Write(buf, order, ph.ReservedEnd); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalPacket is a function that satisfies the ProtocolComponent interface.
// It takes an io.Reader and pulls unmarshals the packet in to the
// ProtocolHeader struct fields. It uses the order parameter to correctly
// unpack the values.
func (ph *ProtocolHeader) UnmarshalPacket(data io.Reader, order binary.ByteOrder) (err error) {
	if err = binary.Read(data, order, &ph.Reserved); err != nil {
		return
	}

	if err = binary.Read(data, order, &ph.Type); err != nil {
		return
	}

	if err = binary.Read(data, order, &ph.ReservedEnd); err != nil {
		return
	}

	return
}
