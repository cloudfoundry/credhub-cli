package models

import (
	"encoding/json"
	"fmt"
)

const JSON_PRETTY_PRINT_INDENT_STRING = "\t"

type Printable interface {
	Terminal() string
	Json() string
}

func prettyPrintJson(m map[string]interface{}) string {
	s, _ := json.MarshalIndent(m, "", JSON_PRETTY_PRINT_INDENT_STRING)
	return string(s)
}

func Println(p Printable, asJson bool) {
	if asJson {
		fmt.Println(p.Json())
	} else {
		fmt.Println(p.Terminal())
	}
}
