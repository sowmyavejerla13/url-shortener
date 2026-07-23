package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/sowmyavejerla13/url-shortener/internal/handler"
	"github.com/sowmyavejerla13/url-shortener/internal/handler/mocks"
)

func TestRegister(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		requestBody string

		serviceErr error

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",

			requestBody: `{
				"name":"John",
				"email":"john@test.com",
				"password":"password123"
			}`,

			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"message":"User registered successfully"
			}`,
		},

		{
			name: "Invalid JSON",

			requestBody: `{
				"name":"John",
				"email":
			}`,

			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"error":"Invalid request body"
			}`,
		},

		{
			name: "Validation Error",

			requestBody: `{
				"name":"",
				"email":"invalid-email",
				"password":"123"
			}`,

			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Service Error",

			requestBody: `{
				"name":"John",
				"email":"john@test.com",
				"password":"password123"
			}`,

			serviceErr: errors.New("email already exists"),

			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"error":"email already exists"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.UserServiceMock{}

			mockService.RegisterFunc = func(name, email, password string) error {
				return tt.serviceErr
			}

			authHandler := handler.NewAuthHandler(mockService)

			req, _ := http.NewRequest(
				http.MethodPost,
				"/register",
				bytes.NewBuffer([]byte(tt.requestBody)),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)

			c.Request = req

			authHandler.Register(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.name == "Validation Error" {

				assert.Contains(t, rec.Body.String(), "errors")

			} else {

				assert.JSONEq(
					t,
					tt.expectedBody,
					rec.Body.String(),
				)

			}
		})
	}
}

func TestLogin(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		requestBody string

		token      string
		serviceErr error

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",

			requestBody: `{
				"email":"john@test.com",
				"password":"password123"
			}`,

			token: "jwt-token",

			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"token":"jwt-token"
			}`,
		},

		{
			name: "Invalid JSON",

			requestBody: `{
				"email":
			}`,

			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"error":"Invalid request body"
			}`,
		},

		{
			name: "Validation Error",

			requestBody: `{
				"email":"invalid-email",
				"password":""
			}`,

			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Invalid Credentials",

			requestBody: `{
				"email":"john@test.com",
				"password":"wrong"
			}`,

			serviceErr: errors.New("invalid email or password"),

			expectedStatus: http.StatusUnauthorized,
			expectedBody: `{
				"error":"invalid email or password"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.UserServiceMock{}

			mockService.LoginFunc = func(email, password string) (string, error) {
				return tt.token, tt.serviceErr
			}

			authHandler := handler.NewAuthHandler(mockService)

			req, _ := http.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBuffer([]byte(tt.requestBody)),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)

			c.Request = req

			authHandler.Login(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.name == "Validation Error" {

				assert.Contains(t, rec.Body.String(), "errors")

			} else {

				assert.JSONEq(
					t,
					tt.expectedBody,
					rec.Body.String(),
				)
			}
		})
	}
}
