package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocastsian/roham/types"
	httpHandler "github.com/gocastsian/roham/userapp/delivery/http"
	"github.com/gocastsian/roham/userapp/mockRepo/userMock"
	"github.com/gocastsian/roham/userapp/service/guard"
	"github.com/gocastsian/roham/userapp/service/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIntegrationHandler_RegisterUser(t *testing.T) {
	type testCase struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}
	testCases := []testCase{
		{
			name: "successfully register a user",
			requestBody: user.RegisterRequest{
				Username:        "testuser1",
				FirstName:       "testuser1",
				LastName:        "testuser1",
				Email:           "testuser1@email.com",
				Password:        "testUser1",
				ConfirmPassword: "testUser1",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"user registered successfully","response":{"user_id":2}}`,
		},
		{
			name: "missing required fields",
			requestBody: user.RegisterRequest{
				Username:        "",
				FirstName:       "",
				LastName:        "",
				Email:           "",
				Password:        "",
				ConfirmPassword: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{"email":"email can not be empty","firstName":"first name can not be empty","lastName":"last name can not be empty","password":"password can not be empty","username":"username can not be empty"},"message":"user validation has error"}`,
		},
		{
			name: "passwords do not match",
			requestBody: user.RegisterRequest{
				Username:        "testuser2",
				FirstName:       "Test",
				LastName:        "User",
				Email:           "testuser2@email.com",
				Password:        "password1",
				ConfirmPassword: "password2",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{"confirm_password":"passwords don't match","password":"password must contain at least one uppercase letter, one lowercase letter, and one number"},"message":"user validation has error"}`,
		},
		{
			name: "invalid email format",
			requestBody: user.RegisterRequest{
				Username:        "testuser3",
				FirstName:       "Test",
				LastName:        "User",
				Email:           "invalid-email",
				Password:        "ValidPass123",
				ConfirmPassword: "ValidPass123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{"email":"must be a valid email address"},"message":"user validation has error"}`,
		},
		{
			name: "username already exists",
			requestBody: user.RegisterRequest{
				Username:        "testuser1",
				FirstName:       "Test",
				LastName:        "User",
				Email:           "unique@email.com",
				Password:        "ValidPass123",
				ConfirmPassword: "ValidPass123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{},"message":"user already exist!"}`,
		},
		{
			name: "email already exists",
			requestBody: user.RegisterRequest{
				Username:        "uniqueuser",
				FirstName:       "Test",
				LastName:        "User",
				Email:           "default@email.com",
				Password:        "ValidPass123",
				ConfirmPassword: "ValidPass123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{},"message":"user already exist!"}`,
		},
		{
			name: "short password",
			requestBody: user.RegisterRequest{
				Username:        "testuser4",
				FirstName:       "Test",
				LastName:        "User",
				Email:           "testuser4@email.com",
				Password:        "123",
				ConfirmPassword: "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":{"password":"the length must be no less than 8"},"message":"user validation has error"}`,
		},
	}
	userRepo := userMock.NewUserRepoMock()
	userValidator := user.NewValidator(userRepo)
	userConf := user.Config{}
	userSrv := user.NewService(userRepo, userValidator, nil, nil, userConf)
	guardSrv := guard.Service{}
	userHandler := httpHandler.NewHandler(userSrv, guardSrv)
	e := echo.New()
	var userRepoCtx context.Context
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			err := userHandler.RegisterUser(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if rec.Code == tc.expectedStatus {
				assert.Equal(t, tc.expectedBody, strings.Trim(rec.Body.String(), "\n"))
				if tc.expectedStatus == http.StatusOK || tc.expectedStatus == http.StatusCreated {
					responseBody := struct {
						Message string   `json:"message"`
						UserID  types.ID `json:"user_id"`
					}{}
					err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
					if err != nil {
						t.Fatal("json unmarshal err for body response: ", err)
					}

					createdUser, err := userRepo.GetUser(userRepoCtx, 2)
					fmt.Printf("response%v", createdUser)
					if err != nil {
						t.Fatal("register user has repository error", err)
					}
					assert.Equal(t, createdUser.Username, tc.requestBody.(user.RegisterRequest).Username)
				}

			}
		})
	}

}
