package eofutil

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
	for i := 0; i < 10; i++ {
		ch <- struct{}{}
	}
	//timer := time.NewTimer(500 * time.Millisecond)
	//select {
	//case :
	//	return
	//case <-timer.C:
	//	return
	//}
}
