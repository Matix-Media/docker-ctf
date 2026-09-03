//go:build linux

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// stdinIsTerminal prueft per ioctl(TCGETS), ob wirklich ein Terminal an stdin haengt.
//
// Warum nicht einfach os.ModeCharDevice? Weil 'docker exec' ohne -i /dev/null an
// stdin haengt — und /dev/null IST ein Character Device. Ohne diese Pruefung
// wuerde der Teilnehmer erst "Druecke ENTER" lesen und danach trotzdem eine
// Fehlermeldung bekommen. Genau dieser Level lebt aber davon, dass der
// Unterschied zwischen 'docker exec' und 'docker exec -it' klar wird.
func stdinIsTerminal() bool {
	var termios [64]byte
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&termios)),
		0, 0, 0,
	)
	return errno == 0
}
