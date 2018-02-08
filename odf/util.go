package odf

import (
	"encoding/json"
	"os"
)

// ReadJSON reads the JSON data from the provided file and returns a
// generic representation of it that can be consumed by templates.
func ReadJSON(fname string) (map[string]interface{}, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	if err := json.NewDecoder(f).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
