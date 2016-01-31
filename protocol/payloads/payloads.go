// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

// Package lifxpayloads is used for marshaling and unmarshaling different LIFX
// protocol payloads to and from the wire, respectively. This package is not
// meant to be consumed by those wanting to interface with their LIFX devices
// in Golang. This package is designed to be used by the LIFX Golang library
// for communicating with devices. Users are meant to consume that package.
//
// At the time of writing, the main LIFX Go package does not exist. This
// package is a prerequisite for the client package.
package lifxpayloads
