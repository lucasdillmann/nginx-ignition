package eventstream

import "github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"

type consistencyProblem struct {
	Message   *i18n.Message `json:"message"`
	FieldPath string        `json:"fieldPath"`
}

type validationEvent struct {
	Message             *i18n.Message        `json:"message"`
	ConsistencyProblems []consistencyProblem `json:"consistencyProblems"`
}

type errorEvent struct {
	Message string `json:"message"`
}
