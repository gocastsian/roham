package user

import (
	"github.com/gocastsian/roham/types"
	"time"
)

type User struct {
	ID           types.ID   `json:"id"`
	Username     string     `json:"username"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	PhoneNumber  string     `json:"phone_number"`
	Email        string     `json:"email"`
	Avatar       string     `json:"avatar"`
	BirthDate    string     `json:"birth_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	IsActive     bool       `json:"is_active"`
	Role         types.Role `json:"role"`
	PasswordHash string     `json:"password_hash"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	ID          types.ID `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Tokens      Tokens   `json:"tokens"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RegisterResponse struct {
	ID types.ID `json:"user_id"`
}
type RegisterRequest struct {
	Username        string `json:"username"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
