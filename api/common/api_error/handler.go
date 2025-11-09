package api_error

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type problemDetail struct {
	FieldPath string `json:"fieldPath"`
	Message   string `json:"message"`
}

func Handler(ctx *gin.Context, outcome any) {
	err, isErr := outcome.(error)
	if !isErr {
		err = errors.New(fmt.Sprintf("%s", outcome))
	}

	httpError := &ApiError{}
	consistencyError := &validation.ConsistencyError{}
	coreError := &core_error.CoreError{}

	switch {
	case errors.As(err, &httpError):
		handleHttpError(ctx, httpError)
	case errors.As(err, &consistencyError):
		handleConsistencyError(ctx, consistencyError)
	case errors.As(err, &coreError):
		handleCoreError(ctx, coreError)
	case errors.Is(err, jwt.ErrSignatureInvalid):
		handleInvalidTokenError(ctx)
	default:
		handleGenericError(ctx, err)
	}
}

func CanHandle(err error) bool {
	if err == nil {
		return false
	}

	httpError := &ApiError{}
	consistencyError := &validation.ConsistencyError{}
	validationError := &validator.ValidationErrors{}
	coreError := &core_error.CoreError{}

	switch {
	case errors.As(err, &httpError),
		errors.As(err, &consistencyError),
		errors.As(err, &validationError),
		errors.As(err, &coreError),
		errors.Is(err, jwt.ErrSignatureInvalid):
		return true
	default:
		return false
	}
}

func handleGenericError(ctx *gin.Context, err error) {
	stack := stacktrace()
	log.Errorf("Error detected while processing request: %s\n%s", err, stack)
	ctx.Status(http.StatusInternalServerError)
}

func handleInvalidTokenError(ctx *gin.Context) {
	ctx.Status(http.StatusUnauthorized)
}

func handleHttpError(ctx *gin.Context, err *ApiError) {
	ctx.JSON(err.StatusCode, gin.H{"message": err.Message})
}

func handleCoreError(ctx *gin.Context, err *core_error.CoreError) {
	if err.UserRelated {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Message})
	} else {
		handleGenericError(ctx, err)
	}
}

func handleConsistencyError(ctx *gin.Context, err *validation.ConsistencyError) {
	details := make([]*problemDetail, len(err.Violations))
	for index, detail := range err.Violations {
		details[index] = &problemDetail{
			FieldPath: detail.Path,
			Message:   detail.Message,
		}
	}

	sendError(ctx, details)
}

func sendError(ctx *gin.Context, details []*problemDetail) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message":             "One or more consistency problems were found",
		"consistencyProblems": details,
	})
}
