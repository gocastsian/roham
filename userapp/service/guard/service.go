package guard

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gocastsian/roham/userapp/service/guard/opa"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Service struct {
	config Config
	logger *slog.Logger
}

func NewAuthService(cfg Config, logger *slog.Logger) Service {
	return Service{
		config: cfg,
		logger: logger,
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

		return []byte("roham"), nil
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

func (srv Service) GuardMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		claims, err := srv.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		input := map[string]interface{}{
			"role":    claims.UserClaim.Role,
			"request": c.QueryParam("request"),
		}

		// check with opa
		if err := opa.PolicyEvaluation(c.Request().Context(), opa.RegoAuthorization, opa.RuleCheckRequestOnly, input); err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
		}

		c.Set(UserCtxKey, claims)
		return next(c)
	}
}
