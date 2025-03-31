package user

import (
	"time"

	"github.com/gocastsian/roham/types"
)

type GetAllUsersItem struct {
	ID          types.ID   `json:"id"`
	Username    string     `json:"username"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	Email       string     `json:"email"`
	Avatar      string     `json:"avatar"`
	BirthDate   string     `json:"birth_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Role        types.Role `json:"role"`
}

type GetAllUsersResponse struct {
	Users []GetAllUsersItem `json:"users"`
}
