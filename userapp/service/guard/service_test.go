package guard_test

import (
	"context"
	"testing"
	"time"

	"github.com/gocastsian/roham/pkg/opa"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/guard"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
	config := guard.Config{
		SignKey:              "testkey",
		AccessSubject:        "access",
		AccessExpirationTime: time.Hour,
	}
	service := guard.NewService(config, nil, nil)
	userClaim := guard.UserClaim{ID: 123, Role: 1}

	token, err := service.CreateAccessToken(userClaim)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestCreateRefreshToken(t *testing.T) {
	config := guard.Config{
		SignKey:              "testkey",
		AccessSubject:        "access",
		AccessExpirationTime: time.Hour,
	}
	service := guard.NewService(config, nil, nil)
	userClaim := guard.UserClaim{ID: 123, Role: 1}

	token, err := service.CreateRefreshToken(userClaim)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	config := guard.Config{
		SignKey: "roham",
	}

	service := guard.NewService(config, nil, nil)

	claims := guard.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "access",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserClaim: guard.UserClaim{ID: 123},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("roham"))

	parsedClaims, err := service.ParseToken("Bearer " + tokenString)
	assert.NoError(t, err)
	assert.Equal(t, types.ID(123), parsedClaims.UserClaim.ID)
}

func TestCheckPolicy(t *testing.T) {
	policy := `package test

      default allow = false

      # Define role constants
      role_admin := 1

      # Admin can access everything
      allow if {
          input.role == role_admin
      }`

	cfg := opa.Config{
		Package: "test",
		Rule:    "allow",
		Policy:  policy,
		IsPath:  false,
	}

	evaluator, err := opa.NewOPAEvaluator(cfg)
	assert.NoError(t, err)

	service := guard.NewService(guard.Config{}, nil, evaluator)

	ctx := context.Background()
	input := map[string]interface{}{"role": 1}

	err = service.CheckPolicy(ctx, input)
	assert.NoError(t, err)
}

func TestCheckPolicy_Error(t *testing.T) {
	policy := `package test

      default allow = false

      # Define role constants
      role_admin := 1

      # Admin can access everything
      allow if {
          input.role == role_admin
    }`

	cfg := opa.Config{
		Package: "test",
		Rule:    "allow",
		Policy:  policy,
		IsPath:  false,
	}

	evaluator, err := opa.NewOPAEvaluator(cfg)
	assert.NoError(t, err)

	service := guard.NewService(guard.Config{}, nil, evaluator)

	ctx := context.Background()
	input := map[string]interface{}{"role": 0}

	err = service.CheckPolicy(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "bindings results[[{[true] map[x:false]}]] ok[true]")
}
