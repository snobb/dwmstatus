package main

// Rewrite of my C version of statusbar.

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"unicode"

	"dwmstatus/pkg/alsa"
	"dwmstatus/pkg/x11"
)

var (
	batPath  string
	wifiPath string
	laPath   string = "/proc/loadavg"
	sprites         = []rune{
		'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█',
	}
)

func loadAverage() string {
	loadavg, err := ioutil.ReadFile(laPath)
	if err != nil {
		log.Println(err)
		return "? ? ?"
	}

	return strings.Join(strings.Fields(string(loadavg))[:3], " ")
}

func wifi() string {
	buf, err := ioutil.ReadFile(wifiPath)
	if err != nil {
		log.Println(err)
		return "down"
	}

	return strings.TrimSpace(string(buf))
}

func batteryStatus() rune {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/status", batPath))
	if err != nil {
		log.Println(err)
		return '?'
	}

	switch unicode.ToLower(rune(buf[0])) {
	case 'c':
		return '+'
	case 'd':
		return '-'
	case 'i':
		return '='
	case 'f':
		return '='
	default:
		return '?'
	}
}

func battery() float64 {
	strnow, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_now", batPath))
	if err != nil {
		log.Println("energy_now", err)
		return -1
	}

	strfull, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_full", batPath))
	if err != nil {
		log.Println("energy_full", err)
		return -1
	}

	var now, full int
	fmt.Sscanf(string(strnow), "%d", &now)
	fmt.Sscanf(string(strfull), "%d", &full)

	return float64(now * 100 / full)
}

func volume() rune {
	vol := alsa.GetVolume()
	if vol < 0 {
		return 'M'
	}

	return sprites[(vol*len(sprites))/100]
}

func main() {
	x11.OpenDisplay()
	defer x11.CloseDisplay()

	var tmpls []string
	var vals []interface{}

	addField := func(tmpl string, val ...interface{}) {
		tmpls = append(tmpls, tmpl)
		vals = append(vals, val...)
	}

	for range time.Tick(time.Second) {
		addField("%s", loadAverage())

		if wifiPath != "" {
			addField("wifi:%s", wifi())
		}

		if batPath != "" {
			addField("bat:%c%0.1f%%", batteryStatus(), battery())
		}

		addField("vol:%c", volume())
		addField("%s", time.Now().Format("Mon Jan 2 15:04:05"))

		x11.SetRootTitle(strings.Join(tmpls, " | "), vals...)
		tmpls, vals = tmpls[:0], vals[:0] // clean line fields
	}
}
