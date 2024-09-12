package internal

import (
    "github.com/iancoleman/strcase"
    "strings"
)

func ToUpper(input string) string {
	return strings.ToUpper(input)
}

func ToCamel(input string) string {
	return strcase.ToCamel(input)
}