package main

import (
	"log"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

func setupLogging() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if len(AppConfig.ProgramLog.LogFile) > 0 {
		log.SetOutput(&lumberjack.Logger{
			Filename:   AppConfig.ProgramLog.LogFile,
			MaxSize:    AppConfig.ProgramLog.LogLimitInMB, // megabytes
			MaxBackups: 5,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	log.Println("Logger Started. OS is " + runtime.GOOS)
}
