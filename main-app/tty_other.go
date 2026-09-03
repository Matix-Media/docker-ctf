//go:build !linux

package main

import "os"

// stdinIsTerminal ist der Fallback fuer lokale Builds ausserhalb von Linux.
// Die Container laufen immer unter Linux und benutzen die ioctl-Variante.
func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	if devNull, err := os.Stat(os.DevNull); err == nil && os.SameFile(info, devNull) {
		return false
	}
	return true
}
