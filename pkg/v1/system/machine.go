package system

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/rayaman/go-qemu/pkg/v1/types"
	"github.com/rayaman/go-qemu/pkg/v1/types/accel"
	"github.com/rayaman/go-qemu/pkg/v1/types/arch"
	"github.com/rayaman/go-qemu/pkg/v1/types/boot"
	"github.com/rayaman/go-qemu/pkg/v1/types/chip"
	"github.com/rayaman/go-qemu/pkg/v1/types/memory"
	"github.com/rayaman/go-qemu/pkg/v1/types/nic"
	"github.com/rayaman/go-qemu/pkg/v1/types/numa"
	"github.com/rayaman/go-qemu/pkg/v1/types/smp"
	"github.com/rayaman/go-qemu/pkg/v1/types/utils"
)

/*
TODO:
-add-fd fd=fd,set=set[,opaque=opaque]

	Add 'fd' to fd 'set'

-set group.id.arg=value

	set <arg> parameter for item <id> of type <group>
	i.e. -set drive.$id.file=/path/to/image

-global driver.property=value
-global driver=driver,property=property,value=value

	set a global default for a driver property
*/
type Machine struct {
	Arch arch.System // The binary we use
	// Amount of memory in MB
	Memory memory.Memory `json:"m,omitempty"`
	// Number of CPU cores
	Cores          smp.SMP     `json:"smp,omitempty"`
	Cpu            chip.CHIP   `json:"cpu,omitempty"`
	Accel          accel.Accel `json:"accel,omitempty"`
	Boot           boot.Boot   `json:"boot,omitempty"`
	Numa           numa.Numa   `json:"numa,omitempty" omit:"true"`
	MemoryPath     string      `json:"memory-path,omitempty"`
	MemoryPrealloc types.Flag  `json:"memory-prealloc,omitempty"`
	Nic            nic.NIC     `json:"nic,omitempty"`

	// Graphics
	NoGraphic types.Flag `json:"nographic,omitempty"`

	// Block devices
	HardDiskA string `json:"hda,omitempty"`
	HardDiskB string `json:"hdb,omitempty"`
	HardDiskC string `json:"hdc,omitempty"`
	HardDiskD string `json:"hdd,omitempty"`
	FloppyA   string `json:"fda,omitempty"`
	FloppyB   string `json:"fdb,omitempty"`
	CDROM     string `json:"cdrom,omitempty"`
}

func (m *Machine) Expand() []string {
	fields := structs.Fields(m)
	exp := []string{}
	for _, field := range fields {
		tag := strings.ReplaceAll(field.Tag("json"), ",omitempty", "")
		omit := field.Tag("omit")
		if tag != "" {
			if field.Kind() == reflect.Struct || field.Kind() == reflect.Interface && !field.IsZero() {
				if omit != "" {
					exp = append(exp, utils.Expand(field.Value())...)
				} else {
					exp = append(exp, utils.Expand(field.Value(), tag)...)
				}
			} else {
				if !field.IsZero() {
					if fmt.Sprintf("%v", field.Value()) == "flag-on" {
						exp = append(exp, "-"+tag)
					} else {
						exp = append(exp, "-"+tag, fmt.Sprintf("%v", field.Value()))
					}
				}
			}
		}
	}
	exp = utils.Remove(exp, "")
	return exp
}
