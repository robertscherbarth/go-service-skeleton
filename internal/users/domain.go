package users

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string
	Tag  string
}
