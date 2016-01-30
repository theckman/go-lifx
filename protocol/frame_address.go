package lifxprotocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

// MaxFrameAddressReserved is the max size of the FrameAddress.Reserved
// field. It only uses the top 6 bits so the maximum value is 63.
const MaxFrameAddressReserved = ^uint8(0) >> 2

// FrameAddressByteSize is the number of bytes in a marshaled FrameAddress struct
const FrameAddressByteSize int = 16

// ErrFrameAddressReservedOverflow is the error returned when the FrameAddress.Reserved value is too large.
//
// Also, what is this... Java?
var ErrFrameAddressReservedOverflow = fmt.Errorf("The Reserved field cannot be larger than %d, suggested value is 0", MaxFrameAddressReserved)

// ErrFrameAddressTargetMalformed is the error returned when the Target field is malformed. In other
// words, it does not contain exactly 6 bytes.
var ErrFrameAddressTargetMalformed = errors.New("The Target byte slice is malformed; the slice must contain 6 bytes")

// FrameAddress is a struct that contains information about the following things:
//
// 		* target device address
// 		* flag specifying whether an ack message is required
// 		* flag specifying whether a state response message is required
// 		* message sequence number
type FrameAddress struct {
	// Target is the devide address (MAC address) we are targetting this packet for.
	// As the device address is a MAC address, this byte slice should consist of 6
	// bytes. If we want to target all devices, this slice should either be empty/nil
	// or 6 bytes with a value of 0.
	//
	// The underlying protocol spec defines this as an 8 byte (uint64) value with
	// the right two-most bytes appearing to be used for padding. While the right-most
	// two bytes look to be used for padding purposes, the spec does not explicitly
	// define them as padding. This slice can be 8 bytes in length *ONLY* if the
	// last two bytes (indicies 6 and 7) are zero (0). This is only to retain some
	// compatibility with how the spec is written.
	Target net.HardwareAddr

	// ReservedBlock is reserved space; must all be zero
	// This entire space equals 48 bits
	ReservedBlock [6]uint8

	// Reserved space specified by the protocol definition
	// This uses the low 6 bits
	Reserved uint8

	// AckRequired: acknowledgement message is required
	AckRequired bool

	// ResRequired: response message is required
	ResRequired bool

	// Sequence is a wrap-around message sequence number
	Sequence uint8
}

func NewFrameAddress() *FrameAddress { return &FrameAddress{} }

func (fra *FrameAddress) MarshalPacket(order binary.ByteOrder) ([]byte, error) {
	if fra.Reserved > MaxFrameAddressReserved {
		return nil, ErrFrameAddressReservedOverflow
	}

	var ack, res uint8

	buf := new(bytes.Buffer)

	var u64 uint64

	// if the length of the target slice is 6
	// or if the length of the target slice is 8
	//    and byte 7 == 0 and byte 8 == 0
	if len(fra.Target) == 6 ||
		(len(fra.Target) == 8 && fra.Target[6] == 0 && fra.Target[7] == 0) {
		u64 = targetSliceToUint(fra.Target)
	} else {

	}

	if err := binary.Write(buf, order, u64); err != nil {
		return nil, err
	}

	for _, value := range fra.ReservedBlock {
		if err := binary.Write(buf, order, value); err != nil {
			return nil, err
		}
	}

	if fra.AckRequired {
		ack = 1
	}

	if fra.ResRequired {
		res = 1
	}

	// the next 8 bit chunk consists of multifple fields:
	// Reserved: 6
	// AckRequired: 1
	// ResponseRequired: 1
	u8 := fra.Reserved<<2 |
		ack<<1 | res

	if err := binary.Write(buf, order, u8); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, order, fra.Sequence); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fra *FrameAddress) UnmarshalPacket(data io.Reader, order binary.ByteOrder) error {
	var u64 uint64

	if err := binary.Read(data, order, &u64); err != nil {
		return err
	}

	fra.Target = uintToTargetSlice(u64)

	for i := range fra.ReservedBlock {
		if err := binary.Read(data, order, &fra.ReservedBlock[i]); err != nil {
			return err
		}
	}

	var u8 uint8

	if err := binary.Read(data, order, &u8); err != nil {
		return err
	}

	fra.Reserved = u8 >> 2         // get top 6 bits
	fra.AckRequired = u8>>1&1 == 1 // get 7th bit and eval if it's true
	fra.ResRequired = u8&1 == 1    // get 8th bit and eval if it's true

	return binary.Read(data, order, &fra.Sequence)
}

func targetSliceToUint(target net.HardwareAddr) uint64 {
	return uint64(target[0])<<55 |
		uint64(target[1])<<47 |
		uint64(target[2])<<39 |
		uint64(target[3])<<31 |
		uint64(target[4])<<23 |
		uint64(target[5])<<15
}

func uintToTargetSlice(u64 uint64) net.HardwareAddr {
	hwaddr := make(net.HardwareAddr, 6)
	hwaddr[0] = byte(u64 >> 55)
	hwaddr[1] = byte(u64 >> 47)
	hwaddr[2] = byte(u64 >> 39)
	hwaddr[3] = byte(u64 >> 31)
	hwaddr[4] = byte(u64 >> 23)
	hwaddr[5] = byte(u64 >> 15)
	return hwaddr
}
