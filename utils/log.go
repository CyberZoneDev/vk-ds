package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type CustomLogger struct {
	*logrus.Entry
}
type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339)
	level := strings.ToUpper(entry.Level.String())
	msg := fmt.Sprintf("%s [%s] | %s\n", timestamp, level, entry.Message)
	return []byte(msg), nil
}

func NewLogger(instancename string) *CustomLogger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&CustomFormatter{})

	log := logger.WithFields(logrus.Fields{
		"instance": instancename,
	})

	return &CustomLogger{log}
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.Entry.Infof(format, args...)
}

func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.Entry.Debugf(format, args...)
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}

func (l *CustomLogger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

func (l *CustomLogger) Fatalf(format string, args ...interface{}) {
	l.Entry.Fatalf(format, args...)
}
