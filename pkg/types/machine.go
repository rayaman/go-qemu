package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
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
	Cores          SMP       `json:"smp,omitempty"`
	Cpu            CHIP      `json:"cpu,omitempty"`
	Accel          Accel     `json:"accel,omitempty"`
	Boot           Boot      `json:"boot,omitempty"`
	Numa           numa.Numa `json:"numa,omitempty" omit:"true"`
	MemoryPath     string    `json:"memory-path,omitempty"`
	MemoryPrealloc Flag      `json:"memory-prealloc,omitempty"`
	Nic            NIC       `json:"nic,omitempty"`

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

func (m *Machine) Expand() string {
	fields := structs.Fields(m)
	exp := []string{}
	for _, field := range fields {
		tag := strings.ReplaceAll(field.Tag("json"), ",omitempty", "")
		omit := field.Tag("omit")
		if tag != "" {
			if field.Kind() == reflect.Struct || field.Kind() == reflect.Interface && !field.IsZero() {
				if omit != "" {
					exp = append(exp, Expand(field.Value()))
				} else {
					exp = append(exp, Expand(field.Value(), tag))
				}
			} else {
				if !field.IsZero() {
					if strings.Contains(fmt.Sprintf("%v", field.Value()), " ") {
						exp = append(exp, fmt.Sprintf(`-%v "%v"`, tag, field.Value()))
					} else {
						exp = append(exp, fmt.Sprintf("-%v %v", tag, field.Value()))
					}
				}
			}
		}
	}

	return strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("qemu-system-%v %v", m.Arch, strings.Join(exp, " ")), " flag-on", ""), "  ", " ")
}

func Expand(obj any, tag ...string) string {
	opts := []string{}
	fields := structs.Fields(obj)
	useSpace := false
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
				useSpace = true
				for _, node := range value {
					opts = append(opts, "-numa "+opt+","+Expand(node))
				}
			case []numa.Dist:
				useSpace = true
				for _, dist := range value {
					opts = append(opts, "-numa "+opt+","+Expand(dist))
				}
			case []numa.CPU:
				useSpace = true
				for _, cpu := range value {
					opts = append(opts, "-numa "+opt+","+Expand(cpu))
				}
			case []numa.HMATLB:
				useSpace = true
				for _, hmatlb := range value {
					opts = append(opts, "-numa "+opt+","+Expand(hmatlb))
				}
			case []numa.HMATCache:
				useSpace = true
				for _, hmatcache := range value {
					opts = append(opts, "-numa "+opt+","+Expand(hmatcache))
				}
			case Options:
				opts = append(opts, value.ExpandOptions())
			case []Drives:
				opts = append(opts, fmt.Sprintf("%v=%v", opt, strings.Join(ConvertDrives(value), ",")))
			default:
				if omit == "" {
					if strings.Contains(fmt.Sprintf("%v", value), " ") {
						opts = append(opts, fmt.Sprintf(`%v="%v"`, opt, value))
					} else {
						opts = append(opts, fmt.Sprintf("%v=%v", opt, value))
					}
				} else {
					if strings.Contains(fmt.Sprintf("%v", value), " ") {
						opts = append(opts, fmt.Sprintf(`"%v"`, value))
					} else {
						opts = append(opts, fmt.Sprintf("%v", value))
					}
				}

			}
		}
	}
	otag := ""
	if len(tag) > 0 {
		otag = "-" + tag[0] + " "
	}
	if useSpace {
		return otag + strings.Join(opts, " ")
	}
	return otag + strings.Join(opts, ",")
}
