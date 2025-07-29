package logger

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	log.SetOutput(os.Stdout)
}
