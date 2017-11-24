package main

import (
	l4g "github.com/bsed/log4go"
	"time"
	//"os/signal"
	//"syscall"
	"os"
	"os/signal"
	"syscall"
)

const logname = "example.xml"
var Log l4g.Logger

func loadLog4goConfig() {
	l4g.Debug("Loading configuration")
	_, err := os.Stat(logname)
	if os.IsNotExist(err) {
		return
	}
	l4g.LoadConfiguration(logname)
}

func Close(){
	Log.Close()
}

func main() {
	/*log = l4g.NewDefaultLogger(l4g.INFO)*/
	// Load the configuration (isn't this easy?)
	//log.LoadConfiguration("./example.xml")
	l4g.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
	loadLog4goConfig()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP)
	go func() {
		for {
			<-sig
			loadLog4goConfig()
		}
	}()

	//Log=make(l4g.Logger)
	//Log.LoadConfiguration("example.xml");

	// And now we're ready!
	for {
		l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
		l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Error("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		time.Sleep(1 * time.Second)

		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

		l4g.Info("About that time, eh chaps?")
		l4g.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	}

}
