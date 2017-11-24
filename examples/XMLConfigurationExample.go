package main

import (
	l4g "github.com/bsed/log4go"
	"time"
)

func main() {
	log := make(l4g.Logger)
	// Load the configuration (isn't this easy?)
	log.LoadConfiguration("example.xml")

	// And now we're ready!
	log.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	log.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	log.Info("About that time, eh chaps?")
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	time.Sleep(1 * time.Second)

	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}
