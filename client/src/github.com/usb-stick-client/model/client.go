package model

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	httpClient = http.Client{
		Timeout: time.Duration(60 * time.Second),
	}
	httpClientNoRedirect = http.Client{
		Timeout: time.Duration(60 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse

		},
	}
	ServiceUrl string
	endpoints  = struct {
		bundle func(etag string) string
	}{
		bundle: func(etag string) string {
			return "https://" + ServiceUrl + "/bundle/" + etag
		},
	}
)

type bundleRequest struct {
	etag string
	method string
	data []byte
}

// for the create we should not follow the redirect. Initial request needs to be made with no data, otherwise we are
// limited by the APi gateway 10MB constraint.
func CreateBundle(etag string, data []byte) error {
	e := getPresignedUrlForPut(&bundleRequest{
		etag:etag,
		data: data,
		method:"DOES NOT MATTER",
	})

	return e
}

func GetBundle(etag string) ([]byte, error) {
	return getPresignedUrl(&bundleRequest{
		etag:etag,
		method:"GET",
	})
}

func getPresignedUrl(request *bundleRequest) ([]byte, error) {
	req, e := http.NewRequest(request.method, endpoints.bundle(request.etag), bytes.NewReader(request.data))
	if e != nil {
		return nil, e
	}

	resp, e := httpClient.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	bodyBytes, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("error body")
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return nil, errors.New(fmt.Sprintf("presigned get request bad status\n"))
	}

	return bodyBytes, nil
}

func getPresignedUrlForPut(request *bundleRequest) error {
	req, e := http.NewRequest("PUT", endpoints.bundle(request.etag), nil)
	if e != nil {
		return e
	}

	resp, e := httpClientNoRedirect.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		bodyBytes, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return e
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return errors.New(fmt.Sprintf("unable to get presigned url"))
	}

	presignedUrl := resp.Header.Get("Location")

	req, e = http.NewRequest("PUT", presignedUrl, bytes.NewReader(request.data))
	if e != nil {
		return e
	}

	resp, e = httpClient.Do(req)
	if e != nil {
		return e
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return e
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return errors.New(fmt.Sprintf("unable to upload\n"))
	}

	return nil
}


