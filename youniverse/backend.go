package youniverse

// Client for dbserver/slowdb

import (
	"os"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Backend struct {
	cacheDir string
	baseURLs []string
}

func NewBackend(base []string, dir string) Backend {
	return Backend{
		cacheDir: dir,
		baseURLs: base,
	}
}

func (b *Backend) Get(key string) (data []byte, err error) {
	var resp *http.Response

	if fcacheFile, err := os.Open(b.cacheDir + key); nil == err {
		if cacheData, err := ioutil.ReadAll(fcacheFile); nil == err {
			// Need to do md5 check
			
			return cacheData, nil
		}
	}

	for _, baseURL := range b.baseURLs {
		resp, err = http.Get(baseURL + key)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			err = errors.New(fmt.Sprint("request ", baseURL+key, " failed, interface result stats: ", resp.StatusCode))
			continue
		}

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		err = nil
		break
	}

	return data, err
}
