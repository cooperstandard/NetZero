package routes

import (
	"time"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB          *database.Queries
	TokenSecret string
	AdminKey    string
	Platform    string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token" omitempty:"refresh_token"`
}

type Group struct {
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
}
