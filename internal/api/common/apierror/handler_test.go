package apierror

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_handler(t *testing.T) {
	consistencyError := &validation.ConsistencyError{
		Violations: []validation.ConsistencyViolation{
			{Path: "name", Message: i18n.Static("Name is required")},
			{Path: "port", Message: i18n.Static("Port is invalid")},
		},
	}
	emptyConsistencyError := &validation.ConsistencyError{
		Violations: []validation.ConsistencyViolation{},
	}
	apiError := New(http.StatusTeapot, i18n.Static("API failure"))
	userCoreError := coreerror.New(i18n.Static("User error"), true)
	systemCoreError := coreerror.New(i18n.Static("System error"), false)
	maxBytesError := &http.MaxBytesError{Limit: 1024}

	tests := []struct {
		err            any
		name           string
		expectedBody   string
		expectedStatus int
	}{
		{
			name:           "APIError",
			err:            apiError,
			expectedStatus: http.StatusTeapot,
			expectedBody:   `{"message":"API failure"}`,
		},
		{
			name:           "wrapped APIError",
			err:            fmt.Errorf("request failed: %w", error(apiError)),
			expectedStatus: http.StatusTeapot,
			expectedBody:   `{"message":"API failure"}`,
		},
		{
			name:           "ConsistencyError",
			err:            consistencyError,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[{"fieldPath":"name","message":"Name is required"},{"fieldPath":"port","message":"Port is invalid"}],"message":"api/common/apierror/consistency-problems"}`,
		},
		{
			name:           "wrapped ConsistencyError",
			err:            fmt.Errorf("validation failed: %w", error(consistencyError)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[{"fieldPath":"name","message":"Name is required"},{"fieldPath":"port","message":"Port is invalid"}],"message":"api/common/apierror/consistency-problems"}`,
		},
		{
			name:           "ConsistencyError without violations",
			err:            emptyConsistencyError,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[],"message":"api/common/apierror/consistency-problems"}`,
		},
		{
			name:           "user related CoreError",
			err:            userCoreError,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"User error"}`,
		},
		{
			name:           "wrapped user related CoreError",
			err:            fmt.Errorf("request failed: %w", error(userCoreError)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"User error"}`,
		},
		{
			name:           "system CoreError",
			err:            systemCoreError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "wrapped system CoreError",
			err:            fmt.Errorf("request failed: %w", error(systemCoreError)),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "MaxBytesError",
			err:            maxBytesError,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "wrapped MaxBytesError",
			err:            fmt.Errorf("request failed: %w", maxBytesError),
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "JWT invalid signature",
			err:            jwt.ErrSignatureInvalid,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "wrapped JWT invalid signature",
			err:            fmt.Errorf("authentication failed: %w", jwt.ErrSignatureInvalid),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "generic error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "string outcome",
			err:            "some panic string",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "integer outcome",
			err:            42,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "nil outcome",
			err:            nil,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := performHandlerOutcome(t, test.err)
			assertHandlerResponse(t, recorder, test.expectedStatus, test.expectedBody)
		})
	}
}

func Test_handlerErrorPrecedence(t *testing.T) {
	apiError := New(http.StatusTeapot, i18n.Static("API failure"))
	consistencyError := validation.NewError(nil)
	userCoreError := coreerror.New(i18n.Static("User error"), true)
	systemCoreError := coreerror.New(i18n.Static("System error"), false)
	maxBytesError := &http.MaxBytesError{Limit: 1024}

	tests := []struct {
		err            error
		name           string
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "APIError before ConsistencyError",
			err: errors.Join(
				apiError,
				consistencyError,
				maxBytesError,
				jwt.ErrSignatureInvalid,
			),
			expectedStatus: http.StatusTeapot,
			expectedBody:   `{"message":"API failure"}`,
		},
		{
			name: "ConsistencyError before CoreError",
			err: errors.Join(
				consistencyError,
				userCoreError,
				maxBytesError,
				jwt.ErrSignatureInvalid,
			),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"consistencyProblems":[],"message":"api/common/apierror/consistency-problems"}`,
		},
		{
			name:           "CoreError before MaxBytesError",
			err:            errors.Join(systemCoreError, maxBytesError, jwt.ErrSignatureInvalid),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "MaxBytesError before JWT error",
			err:            errors.Join(maxBytesError, jwt.ErrSignatureInvalid),
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := performHandlerOutcome(t, test.err)
			assertHandlerResponse(t, recorder, test.expectedStatus, test.expectedBody)
		})
	}
}

func Test_handlerAsRecovery(t *testing.T) {
	tests := []struct {
		outcome        any
		name           string
		expectedBody   string
		expectedStatus int
	}{
		{
			name:           "APIError",
			outcome:        New(http.StatusTeapot, i18n.Static("API failure")),
			expectedStatus: http.StatusTeapot,
			expectedBody:   `{"message":"API failure"}`,
		},
		{
			name:           "MaxBytesError",
			outcome:        &http.MaxBytesError{Limit: 1024},
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "JWT invalid signature",
			outcome:        jwt.ErrSignatureInvalid,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine := gin.New()
			engine.Use(gin.CustomRecoveryWithWriter(nil, Handler))
			engine.GET("/", func(*gin.Context) {
				panic(test.outcome)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			engine.ServeHTTP(recorder, request)

			assertHandlerResponse(t, recorder, test.expectedStatus, test.expectedBody)
		})
	}
}

func Test_canHandle(t *testing.T) {
	apiError := New(http.StatusBadRequest, i18n.Static("message"))
	consistencyError := &validation.ConsistencyError{}
	coreError := &coreerror.CoreError{Message: i18n.Static("message")}
	wrappedAPIError := fmt.Errorf("request failed: %w", error(apiError))
	wrappedConsistencyError := fmt.Errorf("validation failed: %w", error(consistencyError))
	wrappedCoreError := fmt.Errorf("request failed: %w", error(coreError))

	tests := []struct {
		err      error
		name     string
		expected bool
	}{
		{
			name:     "APIError",
			err:      apiError,
			expected: true,
		},
		{
			name:     "wrapped APIError",
			err:      wrappedAPIError,
			expected: true,
		},
		{
			name:     "ConsistencyError",
			err:      consistencyError,
			expected: true,
		},
		{
			name:     "wrapped ConsistencyError",
			err:      wrappedConsistencyError,
			expected: true,
		},
		{
			name:     "CoreError",
			err:      coreError,
			expected: true,
		},
		{
			name:     "wrapped CoreError",
			err:      wrappedCoreError,
			expected: true,
		},
		{
			name:     "MaxBytesError",
			err:      &http.MaxBytesError{Limit: 1024},
			expected: true,
		},
		{
			name:     "wrapped MaxBytesError",
			err:      fmt.Errorf("request failed: %w", &http.MaxBytesError{Limit: 1024}),
			expected: true,
		},
		{
			name:     "JWT invalid signature",
			err:      jwt.ErrSignatureInvalid,
			expected: true,
		},
		{
			name:     "wrapped JWT invalid signature",
			err:      fmt.Errorf("authentication failed: %w", jwt.ErrSignatureInvalid),
			expected: true,
		},
		{
			name:     "generic error",
			err:      errors.New("generic"),
			expected: false,
		},
		{
			name:     "wrapped generic error",
			err:      fmt.Errorf("request failed: %w", errors.New("generic")),
			expected: false,
		},
		{
			name:     "nil error",
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

func performHandlerOutcome(t *testing.T, outcome any) *httptest.ResponseRecorder {
	t.Helper()

	engine := gin.New()
	engine.GET("/", func(ginContext *gin.Context) {
		Handler(ginContext, outcome)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	engine.ServeHTTP(recorder, request)

	return recorder
}

func assertHandlerResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expectedStatus int,
	expectedBody string,
) {
	t.Helper()

	assert.Equal(t, expectedStatus, recorder.Code)
	if expectedBody == "" {
		assert.Empty(t, recorder.Body.String())
		return
	}

	assert.JSONEq(t, expectedBody, recorder.Body.String())
}
