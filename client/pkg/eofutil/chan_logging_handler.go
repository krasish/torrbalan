package eofutil

import "time"

//LoggingChanEOFHandler logs a standard message including DestName and attempts to
//send to ch for 500 * time.Millisecond on call to Handle.
type LoggingChanEOFHandler struct {
	logger LoggingEOFHandler
	ch     chan<- struct{}
}

func NewLoggingChanEOFHandler(destName string, ch chan<- struct{}) LoggingChanEOFHandler {
	return LoggingChanEOFHandler{logger: LoggingEOFHandler{destName}, ch: ch}
}

func (l LoggingChanEOFHandler) Handle() {
	l.logger.Handle()
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
