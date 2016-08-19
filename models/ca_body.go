package models

type CaBody struct {
	ContentType string        `json:"type" binding:"required"`
	Value       *CaParameters `json:"value,omitempty" binding:"required"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
}
