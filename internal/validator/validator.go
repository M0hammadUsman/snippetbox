package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldsErr map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldsErr) == 0
}

func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldsErr == nil {
		v.FieldsErr = make(map[string]string)
	}
	if _, exists := v.FieldsErr[key]; !exists {
		v.FieldsErr[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}