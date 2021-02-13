package eofutil

import "time"

type LoggingChanEOFHandler struct {
	l  LoggingEOFHandler
	ch chan<- struct{}
}

func NewLoggingChanEOFHandler(destName string, ch chan<- struct{}) LoggingChanEOFHandler {
	return LoggingChanEOFHandler{l: LoggingEOFHandler{destName}, ch: ch}
}

func (l LoggingChanEOFHandler) Handle() {
	l.Handle()
	TryWrite(l.ch)
}

func TryWrite(ch chan<- struct{}) {
	timer := time.NewTimer(500 * time.Millisecond)
	select {
	case ch <- struct{}{}:
		return
	case <-timer.C:
		return
	}
}
