package image

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/abdfnx/gosh"
	"github.com/rayaman/go-qemu/pkg/types"
)

func getOptions(m map[string]string) string {
	str := []string{}
	for k, v := range m {
		if len(v) > 0 {
			str = append(str, fmt.Sprintf("%v=%v", k, v))
		}
	}
	n_str := strings.Join(str, ",")
	if len(n_str) > 0 {
		return "-o \"" + strings.Join(str, ",") + "\" "
	}
	return ""
}

func getData(m map[string]string, key string) (string, bool) {
	if d, ok := m[key]; ok {
		delete(m, key)
		if len(d) == 0 {
			return "", false
		}
		return d, true
	}
	return "", false
}

func getMap(q any) map[string]string {
	m := map[string]string{}
	v := reflect.ValueOf(q).Elem()
	for j := 0; j < v.NumField(); j++ { // Go through all fields of struct
		if !v.Field(j).IsZero() {
			index := strings.ReplaceAll(v.Type().Field(j).Tag.Get("json"), ",omitempty", "")
			if v.Field(j).Type() == reflect.TypeOf(true) {
				m[index] = types.SW[v.Field(j).Bool()]
			} else {
				m[index] = v.Field(j).String()
			}
		}
	}
	return m
}

type Options struct {
	IsBaseImage bool `json:"is_base_image"`
}

// create [--object OBJECTDEF] [-q] [-f FMT] [-b BACKING_FILE [-F BACKING_FMT]] [-u] [-o OPTIONS] FILENAME [SIZE]

// Creates an image based of the supplied image
func Create(i Image, size types.Size, opts ...Options) error {

	data := getMap(i)

	var imagename string
	var ok bool

	if imagename, ok = getData(data, "image_name"); !(ok && imagename != "") {
		return fmt.Errorf("invalid image data, ImageName is not optional")
	}

	os.MkdirAll(filepath.Join(filepath.Dir(imagename), "/base"), os.ModePerm)
	basepath := filepath.Dir(imagename)
	basename := filepath.Base(imagename)

	if len(opts) > 0 && opts[0].IsBaseImage {
		imagename = filepath.Join(basepath, "base", basename)
	}

	additional := []string{}
	if backing_file, ok := getData(data, "backing_file"); ok {
		additional = append(additional, "-b "+backing_file+" ")
	}

	if backing_fmt, ok := getData(data, "backing_fmt"); ok {
		additional = append(additional, "-F "+backing_fmt+" ")
	}

	format := string(strings.ReplaceAll(strings.ToLower(reflect.TypeOf(i).String()), "*formats.", ""))
	cmd := fmt.Sprintf("qemu-img create -f %s %v%v%v %v", format, strings.Join(additional, ""), getOptions(data), imagename, size)

	// Todo remove this when done testing
	fmt.Println(cmd)
	err, _, errout := gosh.RunOutput(cmd)

	if errout != "" {
		return fmt.Errorf("%v", errout)
	}

	if err != nil {
		fmt.Println(err)
	}

	return nil
}
