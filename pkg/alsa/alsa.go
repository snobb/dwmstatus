package alsa

// #cgo CFLAGS: -pedantic -O3
// #cgo LDFLAGS: -lasound
// #include "alsa.h"
import "C"

// GetVolume returns a volume as percentage (0-100%) or -1 if muted
func GetVolume() int {
	return int(C.get_volume()) // call extern func - see alsa.c
}
