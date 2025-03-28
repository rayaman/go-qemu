package image

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rayaman/go-qemu/pkg/v1/types/utils"
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
				m[index] = utils.SW[v.Field(j).Bool()]
			} else {
				m[index] = v.Field(j).String()
			}
		}
	}
	return m
}
