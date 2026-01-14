package apierror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type problemDetail struct {
	Message   *i18n.Message `json:"message"`
	FieldPath string        `json:"fieldPath"`
}

func Handler(ctx *gin.Context, outcome any) {
	err, isErr := outcome.(error)
	if !isErr {
		err = fmt.Errorf("%s", outcome)
	}

	httpError := &APIError{}
	consistencyError := &validation.ConsistencyError{}
	coreError := &coreerror.CoreError{}

	switch {
	case errors.As(err, &httpError):
		handleHTTPError(ctx, httpError)
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

	httpError := &APIError{}
	consistencyError := &validation.ConsistencyError{}
	coreError := &coreerror.CoreError{}

	return errors.As(err, &httpError) ||
		errors.As(err, &consistencyError) ||
		errors.As(err, &coreError) ||
		errors.Is(err, jwt.ErrSignatureInvalid)
}

func handleGenericError(ctx *gin.Context, err error) {
	stack := stacktrace()
	log.Errorf("Error detected while processing request: %s\n%s", err, stack)
	ctx.Status(http.StatusInternalServerError)
}

func handleInvalidTokenError(ctx *gin.Context) {
	ctx.Status(http.StatusUnauthorized)
}

func handleHTTPError(ctx *gin.Context, err *APIError) {
	ctx.JSON(err.StatusCode, gin.H{"message": err.Message})
}

func handleCoreError(ctx *gin.Context, err *coreerror.CoreError) {
	if err.UserRelated {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Message})
	} else {
		handleGenericError(ctx, err)
	}
}

func handleConsistencyError(ctx *gin.Context, err *validation.ConsistencyError) {
	details := make([]problemDetail, len(err.Violations))
	for index, detail := range err.Violations {
		details[index] = problemDetail{
			FieldPath: detail.Path,
			Message:   detail.Message,
		}
	}

	sendError(ctx, details)
}

func sendError(ctx *gin.Context, details []problemDetail) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message":             i18n.M(ctx, i18n.K.CommonErrorConsistencyProblems),
		"consistencyProblems": details,
	})
}
