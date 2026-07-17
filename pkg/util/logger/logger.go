package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

type logger struct {
	entry *logrus.Entry
}

func (l *logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}
func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}
func (l *logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}
func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}
func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{entry: l.entry.WithField(key, value)}
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{entry: l.entry.WithFields(fields)}
}

func New(level string, output io.Writer) *logger {
	l := logrus.New()

	lvl, _ := logrus.ParseLevel(level)
	l.SetLevel(lvl)

	l.SetOutput(output)

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02.01.2006 15:04:05",
	})

	entry := logrus.NewEntry(l)

	return &logger{entry: entry}
}
