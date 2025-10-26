// Package routes functions and utilities for serving netzero
package routes

import (
	"database/sql"
	"time"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/google/uuid"
)

type APIConfig struct {
	DBConn 			*sql.DB
	DB          *database.Queries
	TokenSecret string
	AdminKey    string
	Platform    string
}

type User struct {
	ID           uuid.UUID `json:"id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitzero"`
	UpdatedAt    time.Time `json:"updated_at,omitzero"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name,omitempty"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type Group struct {
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
}

type UserID struct{}
