package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/rayaman/go-qemu/pkg/v1/types/boot"
	"github.com/rayaman/go-qemu/pkg/v1/types/nic"
	"github.com/rayaman/go-qemu/pkg/v1/types/numa"
)

// converts bool to on/off format
var SW = map[bool]string{
	true:  "on",
	false: "off",
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
			case nic.Options:
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
