package letsencrypt

import (
	"strings"

	acmelog "github.com/go-acme/lego/v4/log"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

var logAdapterInstance acmelog.StdLogger = &logAdapter{}

type logAdapter struct{}

func (l *logAdapter) Fatal(args ...any) {
	l.Fatalln(args...)
}

func (l *logAdapter) Fatalln(args ...any) {
	log.Errorf("%v", args...)
}

func (l *logAdapter) Fatalf(format string, args ...any) {
	log.Errorf(format, args...)
}

func (l *logAdapter) Print(args ...any) {
	l.Println(args...)
}

func (l *logAdapter) Println(args ...any) {
	log.Infof("%v", args...)
}

func (l *logAdapter) Printf(format string, args ...any) {
	switch {
	case strings.HasPrefix(format, "[INFO] "):
		format = strings.Replace(format, "[INFO] ", "", 1)
		log.Infof(format, args...)

	case strings.HasPrefix(format, "[WARN] "):
		format = strings.Replace(format, "[WARN] ", "", 1)
		log.Warnf(format, args...)

	case strings.HasPrefix(format, "[ERROR] "):
		format = strings.Replace(format, "[ERROR] ", "", 1)
		log.Errorf(format, args...)

	case strings.HasPrefix(format, "[FATAL] "):
		format = strings.Replace(format, "[FATAL] ", "", 1)
		log.Errorf(format, args...)

	default:
		log.Infof(format, args...)
	}
}
