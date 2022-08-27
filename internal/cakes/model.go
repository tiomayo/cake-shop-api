package cakes

import (
	"time"
)

type (
	Cake struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Rating      float64    `json:"rating"`
		Image       *string    `json:"image"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}
	ListRequestDto struct {
		Title       string `query:"title"`
		Description string `query:"description"`
		Offset      int    `query:"offset" validate:"omitempty,gte=0"`
		Limit       int    `query:"limit" validate:"omitempty,gte=0"`
	}
	RequestDto struct {
		Title       string  `json:"title" validate:"required"`
		Description string  `json:"description"`
		Rating      float64 `json:"rating" validate:"omitempty,numeric"`
		Image       string  `json:"image" validate:"omitempty,url"`
	}
	UpdateRequestDto struct {
		ID          int      `param:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Rating      *float64 `json:"rating" validate:"omitempty,numeric"`
		Image       string   `json:"image" validate:"omitempty,url"`
	}
)
