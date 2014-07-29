package mesos

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Log struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

// LogConfig is used to configure the destinations for log channels. Any nil
// writers will be interpreted as discard.
type LogConfig struct {
	Prefix string
	Debug  io.Writer
	Info   io.Writer
	Warn   io.Writer
	Error  io.Writer
}

const logFlags = log.Ldate | log.Ltime | log.Lshortfile

func newLogger(writer io.Writer, prefix string) *log.Logger {
	if writer == nil {
		writer = ioutil.Discard
	}
	return log.New(writer, prefix, logFlags)
}

func NewLog(config LogConfig) Log {
	return Log{
		Debug: newLogger(config.Debug, fmt.Sprintf("[D] %s:", config.Prefix)),
		Info:  newLogger(config.Info, fmt.Sprintf("[I] %s:", config.Prefix)),
		Warn:  newLogger(config.Warn, fmt.Sprintf("[W] %s:", config.Prefix)),
		Error: newLogger(config.Error, fmt.Sprintf("[E] %s:", config.Prefix)),
	}
}
