package formats

import (
	"github.com/rayaman/go-qemu/pkg/v1/types"
	"github.com/rayaman/go-qemu/pkg/v1/types/image"
)

func init() {
	types.RegisterType("raw", Raw{})
}

type Raw struct {
	// Name of the image
	ImageName string `json:"image_name"`
	// Preallocation mode (allowed values: off, falloc, full). falloc mode preallocates space for image by calling posix_fallocate(). full mode preallocates space for image by writing data to underlying storage. This data may or may not be zero, depending on the storage location.
	Preallocation image.Preallocation `json:"preallocation"`
}
