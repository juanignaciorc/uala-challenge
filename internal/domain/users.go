package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Followers []uuid.UUID `json:"followers"`
	Follwing  []uuid.UUID `json:"following"`
	Tweets    []Tweet     `json:"tweets"`
}
