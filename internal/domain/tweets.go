package domain

import "github.com/google/uuid"

type Tweet struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message" validate:"max=280"`
}
