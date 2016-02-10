// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package lifxprotocol

import (
	"encoding/binary"
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	order binary.ByteOrder
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) SetUpSuite(c *C) {
	t.order = binary.LittleEndian
}
