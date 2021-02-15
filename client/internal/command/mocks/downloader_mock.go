package mocks

import (
	"github.com/krasish/torrbalan/client/internal/domain/download"
	"github.com/stretchr/testify/mock"
)

type DownloaderMock struct {
	mock.Mock
}

func (d *DownloaderMock) Download(info download.Info) {
	d.Called(info)
}
