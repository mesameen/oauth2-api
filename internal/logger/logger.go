package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var log *zap.Logger

func InitiLogger() error {
	var err error
	log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	return err
}

func Infof(format string, args ...interface{}) {
	log.Info(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	log.Error(fmt.Sprintf(format, args...))
}

func Panicf(format string, args ...interface{}) {
	log.Panic(fmt.Sprintf(format, args...))
}
