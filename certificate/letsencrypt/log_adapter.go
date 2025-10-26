package letsencrypt

import (
	"strings"

	acmelog "github.com/go-acme/lego/v4/log"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

var logAdapterInstance acmelog.StdLogger = &logAdapter{}

type logAdapter struct{}

func (l *logAdapter) Fatal(args ...interface{}) {
	l.Fatalln(args...)
}

func (l *logAdapter) Fatalln(args ...interface{}) {
	log.Fatalf("%v", args...)
}

func (l *logAdapter) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l *logAdapter) Print(args ...interface{}) {
	l.Println(args...)
}

func (l *logAdapter) Println(args ...interface{}) {
	log.Infof("%v", args...)
}

func (l *logAdapter) Printf(format string, args ...interface{}) {
	switch {
	case strings.HasPrefix(format, "[INFO] "):
		format = strings.Replace(format, "[INFO] ", "", 1)
		log.Infof(format, args...)

	case strings.HasPrefix(format, "[WARN] "):
		format = strings.Replace(format, "[WARN] ", "", 1)
		log.Warnf(format, args...)

	default:
		log.Infof(format, args...)
	}
}
