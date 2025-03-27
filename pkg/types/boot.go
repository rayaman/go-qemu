package types

// floppy (a), hard disk (c), CD-ROM (d), network (n)
type Drives string

func ConvertDrives(d []Drives) []string {
	var res []string
	for _, v := range d {
		res = append(res, string(v))
	}
	return res
}

const (
	CDROM    Drives = "d"
	FLOPPY   Drives = "a"
	HARDDISK Drives = "c"
	NETWORK  Drives = "n"
)

type Boot struct {
	// Order of boot
	Order []Drives `json:"order,omitempty"`
	Once  []Drives `json:"once,omitempty"`
	Menu  *bool    `json:"menu,omitempty"`
	// The file's name that would be passed to bios as logo picture, if menu=on
	SplashFile string `json:"splash,omitempty"`
	// The period that splash picture last if menu=on, unit is ms
	SplashTime uint `json:"splash-time,omitempty"`
	// The timeout before guest reboot when boot failed, unit is ms
	RebootTime uint  `json:"reboot-timeout,omitempty"`
	Strict     *bool `json:"strict,omitempty"`
}
