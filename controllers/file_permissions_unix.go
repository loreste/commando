//go:build !windows
// +build !windows

package controllers

import (
	"os"
)

// SetOwnership sets the ownership of a file or folder on Unix-based systems
func SetOwnership(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}
