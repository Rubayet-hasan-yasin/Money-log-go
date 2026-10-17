package domain

import "time"

const (
	AuthProviderLocal  = "local"
	AuthProviderGoogle = "google"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"user_name"`
	Password     *string   `json:"-"`
	Name         string    `json:"name"`
	GoogleID     *string   `json:"google_id,omitempty"`
	AuthProvider string    `json:"auth_provider"` // "local" or "google"
	Avatar       *string   `json:"avatar,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

