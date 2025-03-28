package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/rayaman/go-qemu/pkg/types/accel"
	"github.com/rayaman/go-qemu/pkg/types/boot"
	"github.com/rayaman/go-qemu/pkg/types/numa"
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
	Arch System // The binary we use
	// Amount of memory in MB
	Memory Memory `json:"m,omitempty"`
	// Number of CPU cores
	Cores          SMP         `json:"smp,omitempty"`
	Cpu            CHIP        `json:"cpu,omitempty"`
	Accel          accel.Accel `json:"accel,omitempty"`
	Boot           boot.Boot   `json:"boot,omitempty"`
	Numa           numa.Numa   `json:"numa,omitempty" omit:"true"`
	MemoryPath     string      `json:"memory-path,omitempty"`
	MemoryPrealloc Flag        `json:"memory-prealloc,omitempty"`
	Nic            NIC         `json:"nic,omitempty"`

	// Graphics
	NoGraphic Flag `json:"nographic,omitempty"`

	// Block devices
	HardDiskA string `json:"hda,omitempty"`
	HardDiskB string `json:"hdb,omitempty"`
	HardDiskC string `json:"hdc,omitempty"`
	HardDiskD string `json:"hdd,omitempty"`
	FloppyA   string `json:"fda,omitempty"`
	FloppyB   string `json:"fdb,omitempty"`
	CDROM     string `json:"cdrom,omitempty"`
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
					exp = append(exp, Expand(field.Value())...)
				} else {
					exp = append(exp, Expand(field.Value(), tag)...)
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

	exp = remove(exp, "")
	return exp
}

func Expand(obj any, tag ...string) []string {
	opts := []string{}
	fields := structs.Fields(obj)

	for _, field := range fields {
		opt := strings.ReplaceAll(field.Tag("json"), ",omitempty", "")
		opt = strings.ReplaceAll(opt, "omitempty", "")
		omit := field.Tag("omit")
		if !field.IsZero() || field.Kind() == reflect.Bool || strings.Contains(opt, "id") && field.Kind() != reflect.String {
			switch value := field.Value().(type) {
			case bool:
				opts = append(opts, fmt.Sprintf("%v=%v", opt, SW[value]))
			case *bool:
				opts = append(opts, fmt.Sprintf("%v=%v", opt, SW[*value]))
			case []string:
				opts = append(opts, fmt.Sprintf("%v=%v", opt, strings.Join(value, "")))
			case numa.CPUS:
				if value.Last == 0 {
					opts = append(opts, fmt.Sprintf("%v=%v", opt, value.First))
				} else {
					opts = append(opts, fmt.Sprintf("%v=%v-%v", opt, value.First, value.Last))
				}
			case []numa.Node:
				for _, node := range value {
					opts = append(opts, "-numa", opt)
					opts = append(opts, Expand(node)...)
				}
			case []numa.Dist:
				for _, dist := range value {
					opts = append(opts, "-numa", opt)
					opts = append(opts, Expand(dist)...)
				}
			case []numa.CPU:
				for _, cpu := range value {
					opts = append(opts, "-numa", opt)
					opts = append(opts, Expand(cpu)...)
				}
			case []numa.HMATLB:
				for _, hmatlb := range value {
					opts = append(opts, "-numa", opt)
					opts = append(opts, Expand(hmatlb)...)
				}
			case []numa.HMATCache:
				for _, hmatcache := range value {
					opts = append(opts, "-numa", opt)
					opts = append(opts, Expand(hmatcache)...)
				}
			case Options:
				opts = append(opts, value.ExpandOptions()...)
			case []boot.Drives:
				if omit == "tag" {
					opts = append(opts, fmt.Sprintf("%v", strings.Join(boot.ConvertDrives(value), ",")))
				} else {
					opts = append(opts, fmt.Sprintf("%v=%v", opt, strings.Join(boot.ConvertDrives(value), ",")))
				}
			default:
				if omit == "" {
					opts = append(opts, fmt.Sprintf("%v=%v", opt, value))
				} else if omit == "tag" {
					opts = append(opts, fmt.Sprintf("%v", value))
				}
			}
		}
	}
	if len(tag) > 0 {
		return []string{"-" + tag[0], strings.Join(opts, ",")}
	} else {
		return []string{strings.Join(opts, ",")}
	}
}
