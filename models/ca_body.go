package models

type CaBody struct {
	ContentType string        `json:"type" binding:"required"`
	Ca          *CaParameters `json:"ca,omitempty" binding:"required"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
}
