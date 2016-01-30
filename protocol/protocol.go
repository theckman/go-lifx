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
