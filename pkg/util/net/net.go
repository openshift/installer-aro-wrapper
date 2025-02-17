package net

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"net"
	"syscall"
)

// Dial returns a dialled connection with its send and receive buffer sizes set.
// If sz <= 0, we leave the default size.
func Dial(network, address string, sz int) (net.Conn, error) {
	return (&net.Dialer{
		Control: func(network, address string, rc syscall.RawConn) error {
			if sz <= 0 {
				return nil
			}

			return setBuffers(rc, sz)
		},
	}).Dial(network, address)
}

// read socket(7)
func setBuffers(rc syscall.RawConn, sz int) error {
	var err2 error
	err := rc.Control(func(fd uintptr) {
		err2 = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, sz)
	})
	if err2 != nil {
		return err2
	}
	if err != nil {
		return err
	}

	err = rc.Control(func(fd uintptr) {
		err2 = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, sz)
	})
	if err2 != nil {
		return err2
	}

	return err
}
