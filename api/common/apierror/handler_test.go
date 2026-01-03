package apierror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func init() {
	_ = log.Init()
}

func Test_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		err            any
		name           string
		expectedBody   string
		expectedStatus int
	}{
		{
			name:           "APIError",
			err:            New(http.StatusBadRequest, "Bad Request"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Bad Request"}`,
		},
		{
			name: "ConsistencyError",
			err: &validation.ConsistencyError{
				Violations: []validation.ConsistencyViolation{
					{Path: "field1", Message: "error1"},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[{"fieldPath":"field1","message":"error1"}],"message":"One or more consistency problems were found"}`,
		},
		{
			name: "CoreError (UserRelated)",
			err: &coreerror.CoreError{
				Message:     "User error",
				UserRelated: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"User error"}`,
		},
		{
			name: "CoreError (Not UserRelated)",
			err: &coreerror.CoreError{
				Message:     "System error",
				UserRelated: false,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name:           "JWT Invalid Signature",
			err:            jwt.ErrSignatureInvalid,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Generic Error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name:           "Non-error outcome",
			err:            "some panic string",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/", func(c *gin.Context) {
				Handler(c, tt.err)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			r.ServeHTTP(w, req)

			assert.Equalf(
				t,
				tt.expectedStatus,
				w.Code,
				"for test '%s': expected status %d but got %d. Body: %s",
				tt.name,
				tt.expectedStatus,
				w.Code,
				w.Body.String(),
			)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func Test_CanHandle(t *testing.T) {
	tests := []struct {
		err      error
		name     string
		expected bool
	}{
		{
			name:     "APIError",
			err:      New(http.StatusBadRequest, "msg"),
			expected: true,
		},
		{
			name:     "ConsistencyError",
			err:      &validation.ConsistencyError{},
			expected: true,
		},
		{
			name:     "CoreError",
			err:      &coreerror.CoreError{},
			expected: true,
		},
		{
			name:     "JWT Invalid Signature",
			err:      jwt.ErrSignatureInvalid,
			expected: true,
		},
		{
			name:     "Generic Error",
			err:      errors.New("generic"),
			expected: false,
		},
		{
			name:     "Nil Error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CanHandle(tt.err))
		})
	}
}
