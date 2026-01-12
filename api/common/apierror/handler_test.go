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
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_handler(t *testing.T) {
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
					{
						Path:    "field1",
						Message: i18n.Raw("error1"),
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[{"fieldPath":"field1","message":"error1"}],"message":"One or more consistency problems were found"}`,
		},
		{
			name: "CoreError (UserRelated)",
			err: &coreerror.CoreError{
				Message:     i18n.Raw("User error"),
				UserRelated: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"User error"}`,
		},
		{
			name: "CoreError (Not UserRelated)",
			err: &coreerror.CoreError{
				Message:     i18n.Raw("System error"),
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine := gin.New()
			engine.GET("/", func(ginContext *gin.Context) {
				Handler(ginContext, test.err)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equalf(
				t,
				test.expectedStatus,
				recorder.Code,
				"for test '%s': expected status %d but got %d. Body: %s",
				test.name,
				test.expectedStatus,
				recorder.Code,
				recorder.Body.String(),
			)
			if test.expectedBody != "" {
				assert.JSONEq(t, test.expectedBody, recorder.Body.String())
			}
		})
	}
}

func Test_canHandle(t *testing.T) {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, CanHandle(test.err))
		})
	}
}
