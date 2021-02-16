package eofutil

import (
	"bufio"
	"io"
)

type EOFHandler interface {
	Handle()
}

func WriteServerCheckEOF(writer *bufio.Writer, s string, stopChan chan<- struct{}) error {
	h := NewLoggingChanEOFHandler("server", stopChan)
	return WriteCheckEOF(writer, s, h)
}

//WriteCheckEOF tries to write s to the given *bufio.Writer. It flushes once it succeeds writing.
//If an io.EOF error occurs while writing/flushing the given EOFHandler is called and error is handled as
//caller desired. If another error occurs, it will be immediately returned for processing by caller.
func WriteCheckEOF(writer *bufio.Writer, s string, handler EOFHandler) error {
	if _, err := writer.WriteString(s); err != nil {
		if err == io.EOF {
			handler.Handle()
			return nil
		} else {
			return err
		}
	}
	err := writer.Flush()
	if err == io.EOF {
		handler.Handle()
		return nil
	}
	return err
}

//ReadCheckEOF tries to read till delim from *bufio.Reader. If an io.EOF error occurs while reading
//the given EOFHandler is called and error is handled as caller desired. If another error occurs,
//it will be immediately returned for processing by caller.
func ReadCheckEOF(reader *bufio.Reader, delim byte, handler EOFHandler) (string, error) {
	read, err := reader.ReadString(delim)
	if err != nil {
		if err == io.EOF {
			handler.Handle()
			return "", nil
		}
		return "", err
	}
	return read, nil
}

func ReadServerCheckEOF(reader *bufio.Reader, delim byte, stopChan chan<- struct{}) (string, error) {
	h := NewLoggingChanEOFHandler("server", stopChan)
	return ReadCheckEOF(reader, delim, h)
}
