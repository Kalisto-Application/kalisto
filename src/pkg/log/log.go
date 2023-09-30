package log

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	l *log.Logger
}

func New() *Logger {
	os.MkdirAll(path.Join(xdg.DataHome, "kalisto.logs"), os.ModePerm)
	l := log.New(&lumberjack.Logger{
		Filename:   path.Join(xdg.DataHome, "kalisto.logs/logs.txt"),
		MaxSize:    5, // megabytes
		MaxBackups: 2,
		MaxAge:     280, //days
	}, "", log.LUTC|log.Lshortfile)
	return &Logger{l: l}
}

func (l *Logger) Print(message string) {
	l.l.Println(message)
}

func (l *Logger) Trace(message string) {
	l.l.Println(message)
}

func (l *Logger) Debug(message string) {
	l.l.Println(message)
}

func (l *Logger) Info(message string) {
	l.l.Println(message)
}

func (l *Logger) Warning(message string) {
	l.l.Println(message)
}

func (l *Logger) Error(message string) {
	l.l.Println(message)
}

func (l *Logger) Fatal(message string) {
	l.l.Println(message)
}
