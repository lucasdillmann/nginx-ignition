package api_error

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type problemDetail struct {
	FieldPath string `json:"fieldPath"`
	Message   string `json:"message"`
}

func Handler(context *gin.Context, outcome any) {
	err, isErr := outcome.(error)
	if !isErr {
		err = errors.New(fmt.Sprintf("%s", outcome))
	}

	httpError := &ApiError{}
	consistencyError := &validation.ConsistencyError{}
	validationError := &validator.ValidationErrors{}
	coreError := &core_error.CoreError{}

	switch {
	case errors.As(err, &httpError):
		handleHttpError(context, httpError)
	case errors.As(err, &consistencyError):
		handleConsistencyError(context, consistencyError)
	case errors.As(err, &validationError):
		handleValidationError(context, validationError)
	case errors.As(err, &coreError):
		handleCoreError(context, coreError)
	default:
		handleGenericError(context, err)
	}
}

func handleGenericError(context *gin.Context, err error) {
	stack := stacktrace()
	log.Error("Error detected while processing request: %s\n%s", err, stack)
	context.Status(http.StatusInternalServerError)
}

func handleHttpError(context *gin.Context, err *ApiError) {
	context.JSON(err.StatusCode, gin.H{"message": err.Message})
}

func handleCoreError(context *gin.Context, err *core_error.CoreError) {
	if err.UserRelated {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Message})
	} else {
		handleGenericError(context, err)
	}
}

func handleConsistencyError(context *gin.Context, err *validation.ConsistencyError) {
	var details = make([]problemDetail, len(err.Violations))
	for index, detail := range err.Violations {
		details[index] = problemDetail{
			FieldPath: detail.Path,
			Message:   detail.Message,
		}
	}

	sendError(context, &details)
}

func handleValidationError(context *gin.Context, _ *validator.ValidationErrors) {
	context.Status(http.StatusBadRequest)
}

func sendError(context *gin.Context, details *[]problemDetail) {
	context.JSON(http.StatusBadRequest, gin.H{
		"message":             "One or more consistency problems were found",
		"consistencyProblems": details,
	})
}
