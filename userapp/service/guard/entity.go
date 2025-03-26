package guard

import (
	"time"

	"github.com/gocastsian/roham/types"
	jwt "github.com/golang-jwt/jwt/v4"
)

const (
	UserCtxKey = "userInfo"
)

type UserClaim struct {
	ID   types.ID   `json:"user_id"`
	Role types.Role `json:"role"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserClaim UserClaim
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject         string        `koanf:"access_subject"`
	RefreshSubject        string        `koanf:"refresh_subject"`
}
