package eofutil

import "log"

type LoggingEOFHandler struct {
	DestName string
}

func (l LoggingEOFHandler) Handle() {
	log.Printf("EOF while writing to %s", l.DestName)
}
