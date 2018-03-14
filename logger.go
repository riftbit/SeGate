package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

func initLogger() {
	allowedTimestampsFormat := map[string]int{
		"Mon Jan _2 15:04:05 2006":            1,
		"Mon Jan _2 15:04:05 MST 2006":        1,
		"Mon Jan 02 15:04:05 -0700 2006":      1,
		"02 Jan 06 15:04 MST":                 1,
		"02 Jan 06 15:04 -0700":               1,
		"Monday, 02-Jan-06 15:04:05 MST":      1,
		"Mon, 02 Jan 2006 15:04:05 MST":       1,
		"Mon, 02 Jan 2006 15:04:05 -0700":     1,
		"2006-01-02T15:04:05Z07:00":           1,
		"2006-01-02T15:04:05.999999999Z07:00": 1,
		"3:04PM":                    1,
		"Jan _2 15:04:05":           1,
		"Jan _2 15:04:05.000":       1,
		"Jan _2 15:04:05.000000":    1,
		"Jan _2 15:04:05.000000000": 1,
	}

	_, ok := allowedTimestampsFormat[config.Log.TimestampFormat]
	if ok == false {
		logPrintln("Wrong Timestamp Format value in config!")
		time.Sleep(10 * time.Millisecond)
	}

	var formatter logrus.Formatter

	switch strings.ToLower(config.Log.Formatter) {
	case "text":
		formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: config.Log.DisableColors, TimestampFormat: config.Log.TimestampFormat}
		break
	case "json":
		formatter = &logrus.JSONFormatter{TimestampFormat: config.Log.TimestampFormat}
		break
	default:
		logPrintln("Error Log config formatter")
		time.Sleep(10 * time.Millisecond)
		break
	}

	level, err := logrus.ParseLevel(config.Log.LogLevel)
	if err != nil {
		logPrintln(err)
		time.Sleep(10 * time.Millisecond)
	}

	Logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Level:     level,
	}
}

func logPrintln(args ...interface{}) {
	isTesting := os.Getenv("TESTING")
	if isTesting != "YES" {
		log.Println(args)
	}
}
