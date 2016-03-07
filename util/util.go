// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// Package lifxutil is a helper package that provides utility functionality
// required by the different subpackges of the lifx package. This utility
// functionality includes shared functions, as well as shared constants.
package lifxutil

import "net"

// HardwareAddrToUint64 converts a net.HardwareAddr and returns a
// uint64 based on the LIFX specification:
//
// The target device address is 8 bytes long, when using the 6 byte MAC address
// then left-justify the value and zero-fill the last two bytes.
func HardwareAddrToUint64(target net.HardwareAddr) uint64 {
	return uint64(target[0])<<55 |
		uint64(target[1])<<47 |
		uint64(target[2])<<39 |
		uint64(target[3])<<31 |
		uint64(target[4])<<23 |
		uint64(target[5])<<15
}

// Uint64ToHardwareAddr converts a uint64 value to a net.HardwareAddr
// based on the LIFX specification rules. See the comment for
// HardwareAddrToUint64 for more info .
func Uint64ToHardwareAddr(u64 uint64) net.HardwareAddr {
	hwaddr := make(net.HardwareAddr, 6)

	hwaddr[0] = byte(u64 >> 55)
	hwaddr[1] = byte(u64 >> 47)
	hwaddr[2] = byte(u64 >> 39)
	hwaddr[3] = byte(u64 >> 31)
	hwaddr[4] = byte(u64 >> 23)
	hwaddr[5] = byte(u64 >> 15)

	return hwaddr
}
