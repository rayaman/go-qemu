package image

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/rayaman/go-qemu/pkg/types"
)

// Reperesents a qemu image
type Image interface {
	//GetType() types.Format
}

type holder struct {
	Format string
	Image  any
}

// Loades a saved Image into memory
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

// Saves an image to disk
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
