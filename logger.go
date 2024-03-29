package main

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *Logger
	loggerOnce = sync.Once{}
)

type Logger struct {
	z *zap.Logger
	s *zap.SugaredLogger
}

func NewLogger(development bool, withArgs ...interface{}) *Logger {
	var logger *Logger

	loggerOnce.Do(func() {
		logger = newLogger(development, withArgs...)
	})

	return logger
}

func newLogger(development bool, withArgs ...interface{}) *Logger {
	var z *zap.Logger

	if development {
		z = zap.Must(zap.NewDevelopment(zap.AddCallerSkip(1)))
	} else {
		z = zap.Must(zap.NewProduction(zap.AddCallerSkip(1)))
	}

	s := z.Sugar().With(withArgs...)

	zap.ReplaceGlobals(z)

	logger = &Logger{
		z: z,
		s: s,
	}

	return logger
}

func GetLogger() *Logger {
	return logger
}

func (l *Logger) Sync() error {
	if l.z == nil {
		return nil
	}

	return l.z.Sync()
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Debugw(msg, args...)
		return
	}
	l.s.Debugw(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Infow(msg, args...)
		return
	}
	l.s.Infow(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Warnw(msg, args...)
		return
	}
	l.s.Warnw(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Errorw(msg, args...)
		return
	}
	l.s.Errorw(msg, args...)
}

func (l *Logger) Panic(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Panicw(msg, args...)
		return
	}
	l.s.Panicw(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	if l.s == nil {
		zap.L().Sugar().Fatalw(msg, args...)
		return
	}
	l.s.Fatalw(msg, args...)
}
