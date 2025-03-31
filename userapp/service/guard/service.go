package guard

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gocastsian/roham/pkg/opa"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Service struct {
	config             Config
	logger             *slog.Logger
	opaPolicyEvaluator *opa.OPAEvaluator
}

func NewService(cfg Config, logger *slog.Logger, opaPolicyEvaluator *opa.OPAEvaluator) Service {
	return Service{
		config:             cfg,
		logger:             logger,
		opaPolicyEvaluator: opaPolicyEvaluator,
	}
}

func (srv Service) CreateAccessToken(userClaim UserClaim) (string, error) {
	return srv.createToken(userClaim, srv.config.AccessSubject, srv.config.AccessExpirationTime)
}

func (srv Service) CreateRefreshToken(userClaim UserClaim) (string, error) {
	return srv.createToken(userClaim, srv.config.RefreshSubject, srv.config.RefreshExpirationTime)
}

func (srv Service) createToken(userClaim UserClaim, subject string, expireDuration time.Duration) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Add issued at (iat) to ensure validity
		},
		UserClaim: userClaim,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(srv.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (srv Service) ParseToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return []byte(srv.config.SignKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (srv Service) GetClaimsFromEchoContext(c echo.Context) *Claims {
	return c.Get(UserCtxKey).(*Claims)
}

// CheckPolicy evaluates the policy based on the user's role and request
func (srv Service) CheckPolicy(ctx context.Context, input map[string]interface{}) error {
	return srv.opaPolicyEvaluator.Evaluate(ctx, input)
}
