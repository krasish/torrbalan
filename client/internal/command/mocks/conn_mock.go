package mocks

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type ConnMock struct {
	mock.Mock
}

func (c *ConnMock) Read(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *ConnMock) Write(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *ConnMock) Close() error {
	args := c.Called()
	return args.Error(0)
}

func (c *ConnMock) LocalAddr() net.Addr {
	args := c.Called()
	return args.Get(0).(net.Addr)
}

func (c *ConnMock) RemoteAddr() net.Addr {
	args := c.Called()
	return args.Get(0).(net.Addr)
}

func (c *ConnMock) SetDeadline(t time.Time) error {
	args := c.Called(t)
	return args.Error(0)
}

func (c *ConnMock) SetReadDeadline(t time.Time) error {
	args := c.Called(t)
	return args.Error(0)
}

func (c *ConnMock) SetWriteDeadline(t time.Time) error {
	args := c.Called(t)
	return args.Error(0)
}
