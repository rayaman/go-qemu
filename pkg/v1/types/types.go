package types

import "fmt"

type (
	Flag string
)

var (
	register = map[string]any{}

	On *bool = func(b bool) *bool {
		return &b
	}(true)

	Off *bool = func(b bool) *bool {
		return &b
	}(false)

	Set Flag = "flag-on"
)

func GetTypes() map[string]any {
	return register
}

func RegisterType(t string, i any) error {
	if _, ok := register[t]; !ok {
		register[t] = i
		return nil
	}
	return fmt.Errorf("type already registered")
}
