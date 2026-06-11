package entity

import "time"

type Metadata struct {
	ID        int64     `json:"id"         validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

type Pagination struct {
	Limit  int32 `json:"limit"  validate:"required"`
	Offset int32 `json:"offset" validate:"required"`
	Total  int64 `json:"total"  validate:"required"`
}
