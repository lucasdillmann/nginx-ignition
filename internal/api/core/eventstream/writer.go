package eventstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const writeTimeout = time.Minute

type Writer struct {
	writer gin.ResponseWriter
}

func New(writer gin.ResponseWriter) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) Prepare() error {
	header := w.writer.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache, no-transform")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
	w.writer.WriteHeader(http.StatusOK)

	w.refreshWriteDeadline()
	w.writer.Flush()

	return nil
}

func (w *Writer) Event(name string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.refreshWriteDeadline()

	if _, err := fmt.Fprintf(w.writer, "event: %s\ndata: %s\n\n", name, data); err != nil {
		return err
	}

	w.writer.Flush()

	return nil
}

func (w *Writer) Comment(comment string) error {
	w.refreshWriteDeadline()

	if _, err := fmt.Fprintf(w.writer, ": %s\n\n", comment); err != nil {
		return err
	}

	w.writer.Flush()

	return nil
}

func (w *Writer) refreshWriteDeadline() {
	controller := http.NewResponseController(w.writer)
	_ = controller.SetWriteDeadline(time.Now().Add(writeTimeout))
}
