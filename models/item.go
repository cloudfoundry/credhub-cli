package models

type Item interface {
	String() string
}

type item struct {
}

func NewItem() Item {
	return item{}
}

func (it item) String() string {
	return ""
}
