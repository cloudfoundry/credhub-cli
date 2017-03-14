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
		fmt.Printf(p.Json())
	} else {
		fmt.Printf(p.Terminal())
	}
}
