package main

// Rewrite of my C version of statusbar.
// Inspired by https://github.com/oniichaNj/go-dwmstatus

// #cgo LDFLAGS: -lX11 -lasound
// #include <X11/Xlib.h>
// #include "../include/getvol.h"
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"unicode"
)

var batPath string
var dpy = C.XOpenDisplay(nil)

func batteryStatus(path string) (rune, error) {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/status", path))
	if err != nil {
		return '?', err
	}

	switch unicode.ToLower(rune(buf[0])) {
	case 'c':
		return '+', nil
	case 'd':
		return '-', nil
	case 'i':
		return '=', nil
	case 'f':
		return '=', nil
	default:
		return '?', nil
	}
}

func battery(path string) (float64, error) {
	strnow, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_now", path))
	if err != nil {
		return -1, err
	}

	strfull, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_full", path))
	if err != nil {
		return -1, err
	}

	var now, full int
	fmt.Sscanf(string(strnow), "%d", &now)
	fmt.Sscanf(string(strfull), "%d", &full)

	return float64(now * 100 / full), nil
}

func loadAverage(file string) (string, error) {
	loadavg, err := ioutil.ReadFile(file)
	if err != nil {
		return "Couldn't read loadavg", err
	}

	return strings.Join(strings.Fields(string(loadavg))[:3], " "), nil
}

func setStatus(format string, args ...interface{}) {
	status := fmt.Sprintf(format, args...)
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), C.CString(status))
	C.XSync(dpy, 1)
}

func volume() rune {
	vol := int(C.get_volume())

	if vol < 0 {
		return 'M'
	}

	sprites := []rune{
		'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█',
	}

	return sprites[(vol*7)/100]
}

func main() {
	if dpy == nil {
		log.Fatal("Can't open display")
	}

	for {
		t := time.Now().Format("Jan 2 2006 15:04:05 MST")
		bat, err := battery(batPath)
		if err != nil {
			log.Println(err)
		}

		batStatus, err := batteryStatus(batPath)
		if err != nil {
			log.Println(err)
		}

		la, err := loadAverage("/proc/loadavg")
		if err != nil {
			log.Println(err)
		}

		vol := volume()
		setStatus("%s | vol:%c | %c%0.1f%% | %s", la, vol, batStatus, bat, t)

		time.Sleep(time.Second)
	}
}
