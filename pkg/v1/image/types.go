package image

type (
	Options struct {
		IsBaseImage bool `json:"is_base_image"`
	}

	// Reperesents a qemu image
	Image interface {
		//GetType() types.Format
	}

	holder struct {
		Format string
		Image  any
	}
)
