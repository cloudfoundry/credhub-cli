package models

import (
	"fmt"
)

type Printable interface {
	Terminal() string
	Json() string
}

func Println(p Printable, asJson bool) {
	if asJson {
		fmt.Println(p.Json())
	} else {
		fmt.Println(p.Terminal())
	}
}
