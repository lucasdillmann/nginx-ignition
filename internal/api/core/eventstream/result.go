package eventstream

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
)

type streamResult struct {
	payload any
	err     error
	panic   any
}

func drainEvents(writer *Writer, events <-chan progress.Event) error {
	for {
		select {
		case event := <-events:
			if err := writer.Event("log", event); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func writeTerminal(ctx context.Context, writer *Writer, result streamResult) {
	var consistencyError *validation.ConsistencyError
	if errors.As(result.err, &consistencyError) {
		problems := make([]consistencyProblem, 0, len(consistencyError.Violations))
		for _, violation := range consistencyError.Violations {
			problems = append(problems, consistencyProblem{
				Message:   violation.Message,
				FieldPath: violation.Path,
			})
		}

		_ = writer.Event("validation", validationEvent{
			Message:             i18n.M(ctx, i18n.K.ApiCommonApierrorConsistencyProblems),
			ConsistencyProblems: problems,
		})
		return
	}

	if result.panic != nil {
		_ = writer.Event("error", errorEvent{
			Message: fmt.Sprintf("%v", result.panic),
		})
		return
	}

	if result.payload != nil {
		_ = writer.Event("result", result.payload)
		return
	}

	message := "Unable to complete the operation"
	if result.err != nil {
		message = result.err.Error()
	}

	_ = writer.Event("error", errorEvent{Message: message})
}
