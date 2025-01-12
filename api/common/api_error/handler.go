package api_error

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type problemDetail struct {
	FieldPath string `json:"fieldPath"`
	Message   string `json:"message"`
}

func Handler(context *gin.Context, outcome any) {
	err, isErr := outcome.(error)
	if !isErr || !tryHandle(context, err) {
		stack := stacktrace()
		log.Printf("Error detected while processing request: %s\n%s", err, stack)
		context.Status(http.StatusInternalServerError)
	}
}

func tryHandle(context *gin.Context, err error) bool {
	httpError := &ApiError{}
	consistencyError := &validation.ConsistencyError{}
	validationError := &validator.ValidationErrors{}

	switch {
	case errors.As(err, &httpError):
		handleHttpError(context, httpError)
		return true
	case errors.As(err, consistencyError):
		handleConsistencyError(context, consistencyError)
		return true
	case errors.As(err, validationError):
		handleValidationError(context, validationError)
		return true
	}

	return false
}

func handleHttpError(context *gin.Context, err *ApiError) {
	context.JSON(err.StatusCode, gin.H{"message": err.Message})
}

func handleConsistencyError(context *gin.Context, err *validation.ConsistencyError) {
	var details = make([]problemDetail, len(err.Violations))
	for _, detail := range err.Violations {
		details = append(details, problemDetail{
			FieldPath: detail.Path,
			Message:   detail.Message,
		})
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
