package credhub

type mode string

const (
	Overwrite  mode = "overwrite"
	NoOverwrite mode = "no-overwrite"
)