// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxutil_test

import (
	"net"
	"testing"

	"github.com/theckman/go-lifx/util"

	. "gopkg.in/check.v1"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (*TestSuite) Test_HardwareAddrToUint64(c *C) {
	var u64 uint64

	hwaddr, err := net.ParseMAC("01:23:45:67:89:ab")
	c.Assert(err, IsNil)

	result := uint64(hwaddr[0])<<55 |
		uint64(hwaddr[1])<<47 |
		uint64(hwaddr[2])<<39 |
		uint64(hwaddr[3])<<31 |
		uint64(hwaddr[4])<<23 |
		uint64(hwaddr[5])<<15

	u64 = lifxutil.HardwareAddrToUint64(hwaddr)
	c.Check(u64, Equals, result)
}

func (*TestSuite) Test_Uint64ToHardwareAddr(c *C) {
	var hw, hwaddr net.HardwareAddr

	hw, err := net.ParseMAC("01:23:45:67:89:ab")
	c.Assert(err, IsNil)

	u64 := uint64(hw[0])<<55 |
		uint64(hw[1])<<47 |
		uint64(hw[2])<<39 |
		uint64(hw[3])<<31 |
		uint64(hw[4])<<23 |
		uint64(hw[5])<<15

	hwaddr = lifxutil.Uint64ToHardwareAddr(u64)
	c.Check(hwaddr.String(), Equals, hw.String())
}
