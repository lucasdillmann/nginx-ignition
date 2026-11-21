package progress

import "time"

type Level string

type Event struct {
	Time    time.Time `json:"time"`
	Level   Level     `json:"level"`
	Message string    `json:"message"`
}
