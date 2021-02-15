package mocks

import "github.com/stretchr/testify/mock"

type UploaderMock struct {
	mock.Mock
}

func (u *UploaderMock) AddFile(filePath string) (name string, hash string, err error) {
	args := u.Called(filePath)
	return args.String(0), args.String(1), args.Error(2)
}

func (u *UploaderMock) RemoveFile(fileName string) {
	u.Called(fileName)
}
