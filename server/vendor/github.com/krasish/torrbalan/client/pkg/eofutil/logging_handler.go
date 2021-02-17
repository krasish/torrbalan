package eofutil

import "log"

//LoggingEOFHandler logs a standard message including DestName on call to Handle
type LoggingEOFHandler struct {
	DestName string
}

func (l LoggingEOFHandler) Handle() {
	log.Printf("EOF while writing to %s", l.DestName)
}
