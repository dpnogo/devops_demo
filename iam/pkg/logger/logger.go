package logger

import (
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type LogCfg struct {
	LogLevel string
}

// NewLog 初始化 log 级别
func NewLog(logCfg LogCfg) {
	switch strings.ToLower(logCfg.LogLevel) {
	case "err", "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				file := path.Base(frame.File)
				return "", " " + file + ":" + strconv.Itoa(frame.Line)
			},
		})
	}
}
