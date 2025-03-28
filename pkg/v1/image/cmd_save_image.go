package image

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
)

// Saves an image structure to disk
func SaveImage(path string, img Image) error {

	data, err := json.MarshalIndent(holder{
		Format: strings.ReplaceAll(strings.ToLower(reflect.TypeOf(img).String()), "*formats.", ""),
		Image:  img,
	}, "", "\t")

	if err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
