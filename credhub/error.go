package credhub

// Error provides errors for the CredHub client
type Error struct {
	Message string `json:"error"`
}

func (e *Error) Error() string {
	return e.Message
}
