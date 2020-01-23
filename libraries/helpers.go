package libraries

import (
	"io/ioutil"
	"os"
)

func ReadFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("HOME") + path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
