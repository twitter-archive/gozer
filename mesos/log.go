package mesos

import (
	"io"
	"fmt"
	"log"
)

type Log struct {
	Debug	*log.Logger
	Info	*log.Logger
	Warn	*log.Logger
	Error	*log.Logger
}

const logFlags = log.Ldate | log.Ltime | log.Lshortfile

func NewLog(prefix string, debugWriter, infoWriter, warnWriter, errorWriter io.Writer) *Log {
	return &Log{
		Debug:	log.New(debugWriter, fmt.Sprintf("[D] %s:", prefix), logFlags),
		Info:	log.New(infoWriter, fmt.Sprintf("[I] %s:", prefix), logFlags),
		Warn:	log.New(warnWriter, fmt.Sprintf("[W] %s:", prefix), logFlags),
		Error:	log.New(errorWriter, fmt.Sprintf("[E] %s:", prefix), logFlags),
	}
}

