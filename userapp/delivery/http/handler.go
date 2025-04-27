package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/statuscode"
	"github.com/gocastsian/roham/pkg/validator"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/guard"
	"github.com/gocastsian/roham/userapp/service/user"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	UserService  user.Service
	guardService guard.Service
}

func NewHandler(userSrv user.Service, guardService guard.Service) Handler {
	return Handler{
		UserService:  userSrv,
		guardService: guardService,
	}
}

func (h Handler) GetAllUsers(c echo.Context) error {

	res, err := h.UserService.GetAllUsers(c.Request().Context())
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h Handler) Login(c echo.Context) error {
	var request user.LoginRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorResponse{Message: errmsg.ErrInvalidRequestFormat.Error()})
	}

	resp, err := h.UserService.Login(c.Request().Context(), request)
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) authenticate(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, errmsg.ErrorResponse{Message: errmsg.ErrUnauthorized.Error()})
	}

	claim, err := h.guardService.ParseToken(authHeader)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, errmsg.ErrorResponse{Message: errmsg.ErrUnauthorized.Error()})
	}

	userInfo := guard.UserClaim{
		ID:   claim.UserClaim.ID,
		Role: claim.UserClaim.Role,
	}

	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil
	}

	base64Encoded := base64.StdEncoding.EncodeToString(jsonData)

	c.Response().Header().Set("X-User-Info", base64Encoded)

	return c.NoContent(http.StatusOK)
}

func (h Handler) authorize(c echo.Context) error {
	userInfoHeader := c.Request().Header.Get("X-User-Info")
	if userInfoHeader == "" {
		return c.JSON(http.StatusUnauthorized, errmsg.ErrorResponse{Message: errmsg.ErrUnauthorized.Error()})
	}

	decoded, err := base64.StdEncoding.DecodeString(userInfoHeader)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errmsg.ErrorResponse{Message: errmsg.ErrFailedDecodeBase64.Error()})
	}

	var userClaim guard.UserClaim
	if err := json.Unmarshal(decoded, &userClaim); err != nil {
		return c.JSON(http.StatusUnauthorized, errmsg.ErrorResponse{Message: errmsg.ErrFailedUnmarshalJson.Error()})
	}

	input, err := preparePolicyInput(c, userClaim)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorResponse{Message: errmsg.ErrInvalidRequestFormat.Error()})
	}

	// Evaluate policy
	if err := h.guardService.CheckPolicy(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusForbidden, errmsg.ErrorResponse{Message: errmsg.ErrUnauthorized.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func preparePolicyInput(c echo.Context, userClaim guard.UserClaim) (map[string]interface{}, error) {
	queryParams := make(map[string]interface{})
	for param, values := range c.QueryParams() {
		if len(values) == 1 {
			queryParams[param] = values[0]
		} else {
			queryParams[param] = values
		}
	}

	return map[string]interface{}{
		"user": map[string]interface{}{
			"role": userClaim.Role,
			"id":   userClaim.ID,
		},
		"request": map[string]interface{}{
			"method": c.Request().Method,
			"path":   c.Path(),
			"query":  queryParams,
		},
	}, nil
}

func (h Handler) RegisterUser(c echo.Context) error {
	var request user.RegisterRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorResponse{Message: err.Error()})
	}

	response, err := h.UserService.RegisterUser(c.Request().Context(), request)
	if err != nil {

		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), map[string]interface{}{
			"message": err.Error(),
			"errors":  err.(errmsg.ErrorResponse).Errors,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "user registered successfully",
		"response": response,
	})
}

func (h Handler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ID not found.",
		})
	}

	res, err := h.UserService.GetUser(c.Request().Context(), types.ID(uint(userID)))
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.JSON(http.StatusOK, res)

}

func (h Handler) GetCurrentUser(c echo.Context) error {
	claims := h.guardService.GetClaimsFromEchoContext(c)

	currentUser := user.GetAllUsersItem{ID: claims.UserClaim.ID}

	res, err := h.UserService.GetUser(c.Request().Context(), currentUser.ID)
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.JSON(http.StatusOK, res)
}
