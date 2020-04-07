package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func Get(src string) ([]byte, error) {
	if isHTTPSource(src) {
		return Download(src)
	}

	return Read(src)
}

func isHTTPSource(src string) bool {
	if len([]rune(src)) > 3 && src[:4] == "http" {
		return true
	}
	return false
}

func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	return data, errors.WithStack(err)
}

func Download(url string) ([]byte, error) {
	var client http.Client
	res, err := client.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)
		return data, errors.WithStack(err)
	}

	return nil, errors.New(fmt.Sprintf("bad response code: %d", res.StatusCode))
}
