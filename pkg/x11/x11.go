package x11

// #cgo CFLAGS: -pedantic -O3
// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"

import (
	"fmt"
	"log"
	"unsafe"
)

var dpy *C.Display

// OpenDisplay opens a display
func OpenDisplay() {
	dpy = C.XOpenDisplay(nil)
	if dpy == nil {
		log.Fatal("Can't open display")
	}
}

// CloseDisplay closes the display if opened
func CloseDisplay() {
	if dpy == nil {
		return
	}

	C.XCloseDisplay(dpy)
}

// SetRootTitle sets the title of the root window.
func SetRootTitle(format string, args ...interface{}) {
	status := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(status))

	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), status)
	C.XSync(dpy, 1)
}
