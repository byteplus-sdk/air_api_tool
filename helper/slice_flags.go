package helper

import "fmt"

type StringSliceFlag struct {
	value *[]string
}

func NewStringSliceFlag(value *[]string) *StringSliceFlag {
	return &StringSliceFlag{value: value}
}

func (f *StringSliceFlag) GetValue() []string {
	return *f.value
}

// Will be displayed as the default value when use --help (if it is not an empty string)
func (f *StringSliceFlag) String() string {
	if len(*f.value) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", *f.value)
}

func (f *StringSliceFlag) Set(value string) error {
	*f.value = append(*f.value, value)
	return nil
}

func (f *StringSliceFlag) Type() string {
	return "string array"
}
