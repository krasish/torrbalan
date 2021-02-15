package logutil

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogOnErr(t *testing.T) {
	var (
		errorMessage = "test error message"
	)

	t.Run("logs when error is returned", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log.SetOutput(buf)
		f := func() error { return errors.New(errorMessage) }
		LogOnErr(f)
		assert.Contains(t, buf.String(), errorMessage)
	})
	t.Run("does not log when error is not returned", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log.SetOutput(buf)
		f := func() error { return nil }
		LogOnErr(f)
		assert.Equal(t, "", buf.String())
	})

}
