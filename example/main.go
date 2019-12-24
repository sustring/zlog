package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zjmnssy/zlog"
)

type exitFunc func()

func securityExitProcess(exitFunc exitFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Printf("\n[ INFO ] (system) - security exit by %s signal.\n", s)
			exitFunc()
		default:
			fmt.Printf("\n[ INFO ] (system) - unknow exit by %s signal.\n", s)
			exitFunc()
		}
	}
}

func quit() {
	os.Exit(0)
}

func main() {
	var config zlog.Config
	config.LogFilePath = "/home/nssy/Work/4-zjmnssy/logtest"
	config.LogLevels = 0x3f
	config.LogFilePrefix = "test"
	config.Modules = append(config.Modules, "screen")
	config.Modules = append(config.Modules, "file")
	// config.Modules = append(config.Modules, ModulesAll)
	// config.Modules = append(config.Modules, ModulesNone)
	config.IsTimeSplit = true
	config.SplitPeriod = 10
	config.IsSizeSplit = true
	config.SplitSize = 128
	config.IsClear = true
	config.SavePeriod = 1

	zlog.Init(config)

	zlog.Prints(zlog.Info, "screen", "111 - %d", 1)
	zlog.Prints(zlog.Notice, "file", "222 - %s", "2")
	zlog.Prints(zlog.Warn, "screen", "333 - %d", 3)
	zlog.Prints(zlog.Error, "file", "444 - %d", 4)
	zlog.Prints(zlog.Critical, "screen", "555 - %d", 5)

	zlog.Printf(zlog.Info, "file", "111 - %d %d", 1)
	zlog.Printf(zlog.Notice, "screen", "222 - %d", 2)
	zlog.Printf(zlog.Warn, "file", "333 - %d", 3)
	zlog.Printf(zlog.Error, "screen", "444 - %d", 4)
	zlog.Printf(zlog.Critical, "file", "555 - %d", 5)

	go func() {
		for {
			t := time.Now()
			timeNow := t.UTC().Unix()
			zlog.Printf(zlog.Error, "file", "timeNow = %d", timeNow)
			time.Sleep(2 * time.Second)
		}
	}()

	securityExitProcess(quit)
}
