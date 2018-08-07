package utils

import (
	"encoding/json"
	"io/ioutil"
)

// LoadJsonFile load a json file to obj
func LoadJsonFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
