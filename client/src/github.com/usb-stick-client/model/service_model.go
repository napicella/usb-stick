package model

import (
	"github.com/google/uuid"
	"io/ioutil"
)

type UsbService interface {
	Store(filepath string) (string, error)
	Download(etag string) ([]byte, error)
}

func NewUsbService() UsbService {
	return &defaultUsbService{}
}

type defaultUsbService struct {
}

func (*defaultUsbService) Store(filepath string) (string, error) {
	// instead of generating the ID we could use dynamo db to retrieve a UUID
	randUuid, e := uuid.NewUUID()
	if e != nil {
		return "", e
	}
	etag := randUuid.String()
	data, e := ioutil.ReadFile(filepath)
	if e != nil {
		return "", e
	}

	e = CreateBundle(etag, data)
	if e != nil {
		return "", e
	}

	return etag, nil
}

func (*defaultUsbService) Download(etag string) ([]byte, error) {
	body, e := GetBundle(etag)
	if e != nil {
		return nil, e
	}

	return body, nil
}

