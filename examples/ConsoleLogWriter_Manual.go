package main

import (
	"time"
)

import l4g "github.com/bsed/log4go"

func main() {
	//log := l4g.NewLogger()
	log := make(l4g.Logger)
	defer log.Close()

	log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	time.Sleep(1 * time.Second)

	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	log.Finest("Finest")
	log.Fine("Fine")
	log.Debug("Debug")
	log.Trace("Trace")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Critical("Critical")
	// ConsoleLogWriter 使用带缓存的chan进行输出的管理，所以不进行sleep操作
	//  控制台没有输出，因为我们不知道什么时候channel会进行处理
	time.Sleep(1*time.Second)
	log.Close()

	// This makes sure the filters is running
	// time.Sleep(200 * time.Millisecond)
}
