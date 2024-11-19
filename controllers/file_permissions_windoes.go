//go:build windows
// +build windows

package controllers

import (
	"errors"
)

// SetOwnership is a no-op for Windows since os.Chown is not supported
func SetOwnership(path string, uid, gid int) error {
	return errors.New("os.Chown is not supported on Windows")
}
