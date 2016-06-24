package models

type CaBody struct {
	Ca        *CaParameters `json:"root,omitempty" binding:"required"`
	UpdatedAt string        `json:"updated_at,omitempty"`
}
