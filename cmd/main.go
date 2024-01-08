package main

// Rewrite of my C version of statusbar.

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"dwmstatus/pkg/alsa"
	"dwmstatus/pkg/x11"
)

const FILE_BATTERY_NOW = "charge_now"
const FILE_BATTERY_FULL = "charge_full"
const SUSPEND_TIMEOUT = 60
const SUSPEND_THRESHOLD = 10
const SUSPEND_CMD = "/usr/local/bin/suspend.sh"

var (
	batPath  string
	wifiPath string
	laPath   = "/proc/loadavg"
	sprites  = []rune{
		'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█',
	}
)

func loadAverage() string {
	loadavg, err := os.ReadFile(laPath)
	if err != nil {
		log.Println(err)
		return "? ? ?"
	}

	return strings.Join(strings.Fields(string(loadavg))[:3], " ")
}

func wifi() string {
	buf, err := os.ReadFile(wifiPath)
	if err != nil {
		log.Println(err)
		return "down"
	}

	return strings.TrimSpace(string(buf))
}

func batteryStatus() rune {
	buf, err := os.ReadFile(fmt.Sprintf("%s/status", batPath))
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
	strnow, err := os.ReadFile(fmt.Sprintf("%s/%s", batPath, FILE_BATTERY_NOW))
	if err != nil {
		log.Println("energy_now", err)
		return -1
	}

	strfull, err := os.ReadFile(fmt.Sprintf("%s/%s", batPath, FILE_BATTERY_FULL))
	if err != nil {
		log.Println("energy_full", err)
		return -1
	}

	now, err := strconv.Atoi(string(strnow[:len(strnow)-1]))
	if err != nil {
		log.Println("energy_now:atoi", err)
		return -1
	}

	full, err := strconv.Atoi(string(strfull[:len(strfull)-1]))
	if err != nil {
		log.Println("energy_full:atoi", err)
		return -1
	}

	return float64(now * 100 / full)
}

func volume() rune {
	vol := alsa.GetVolume()
	if vol < 0 {
		return 'M'
	}

	return sprites[(vol*(len(sprites)-1))/100]
}

func main() {
	x11.OpenDisplay()
	defer x11.CloseDisplay()

	var tmpls []string
	var vals []interface{}
	var timer = 0

	addField := func(tmpl string, val ...interface{}) {
		tmpls = append(tmpls, tmpl)
		vals = append(vals, val...)
	}

	for range time.Tick(time.Second) {
		addField("%s", loadAverage())

		addField("vol:%c", volume())

		if wifiPath != "" {
			addField("wifi:%s", wifi())
		}

		if batPath != "" {
			status := batteryStatus()
			charge := battery()

			if status == '-' && charge <= SUSPEND_THRESHOLD {
				addField("LOW BATTERY[%0.1f%%]:suspending in %ds", charge, SUSPEND_TIMEOUT-timer)
				timer++

				if timer > SUSPEND_TIMEOUT {
					_ = exec.Command("/bin/sh", SUSPEND_CMD).Run()
				}
			} else {
				addField("bat:%c%0.1f%%", status, charge)
				timer = 0
			}
		}

		addField("%s", time.Now().Format("Mon Jan 2 15:04:05"))

		x11.SetRootTitle(strings.Join(tmpls, " | "), vals...)
		tmpls, vals = nil, nil
	}
}
