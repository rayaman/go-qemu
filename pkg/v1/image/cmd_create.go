package image

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/abdfnx/gosh"
	"github.com/rayaman/go-qemu/pkg/v1/types/disk"
)

// create [--object OBJECTDEF] [-q] [-f FMT] [-b BACKING_FILE [-F BACKING_FMT]] [-u] [-o OPTIONS] FILENAME [SIZE]

// Creates an image based of the supplied image structure
func Create(i Image, size disk.Size, opts ...Options) error {

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
