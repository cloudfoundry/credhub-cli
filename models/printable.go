package models

type Printable interface {
	Terminal() string
	Json() string
}
