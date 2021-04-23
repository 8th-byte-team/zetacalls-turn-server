package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Log formatter
var Log = makeLog()

func makeLog() *logrus.Logger {
	log := logrus.New()

	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	}

	return log
}

func Pretty(things interface{}) {
	s, _ := json.MarshalIndent(things, "", "\t")
	Log.Info("\n", string(s), "\n")
}
