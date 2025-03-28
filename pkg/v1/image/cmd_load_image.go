package image

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/rayaman/go-qemu/pkg/v1/types"
)

// Loads a saved image structure into memory
func LoadImage(path string) (Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var h = &holder{}

	err = json.Unmarshal(data, h)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(h.Image)
	if err != nil {
		return nil, err
	}

	t := reflect.New(reflect.TypeOf(types.GetTypes()[h.Format])).Interface()
	err = json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
