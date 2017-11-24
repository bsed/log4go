package main

import (
	"time"
)

import l4g "github.com/bsed/log4go"

func main() {
	log := make(l4g.Logger)
	log.AddFilter("network", l4g.FINEST, l4g.NewSocketLogWriter("udp", "172.19.13.222:12124"))

	// Run `nc -u -l -p 12124` or similar before you run this to see the following message
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	// This makes sure the output stream buffer is written
	//log.Close()
}
