package eventstream

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type streamFrame struct {
	Data  map[string]any
	Event string
}

func Test_Serve(t *testing.T) {
	t.Run("streams progress events before the result", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest("POST", "/", nil)

		Serve(ctx, func(runContext context.Context) (any, error) {
			progress.Report(runContext, progress.LevelDebug, "debug message")
			progress.Report(runContext, progress.LevelInfo, "info message")

			return map[string]any{"success": true}, nil
		})

		frames := parseStreamFrames(t, recorder.Body.Bytes())
		require.Len(t, frames, 3)

		assert.Equal(t, "log", frames[0].Event)
		assert.Equal(t, "debug", frames[0].Data["level"])
		assert.Equal(t, "debug message", frames[0].Data["message"])

		assert.Equal(t, "log", frames[1].Event)
		assert.Equal(t, "info", frames[1].Data["level"])
		assert.Equal(t, "info message", frames[1].Data["message"])

		assert.Equal(t, "result", frames[2].Event)
		assert.Equal(t, true, frames[2].Data["success"])
	})

	t.Run("streams validation problems", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest("POST", "/", nil)

		Serve(ctx, func(runContext context.Context) (any, error) {
			return nil, validation.SingleFieldError(
				"domainNames",
				i18n.M(runContext, i18n.K.CommonCannotBeEmpty),
			)
		})

		frames := parseStreamFrames(t, recorder.Body.Bytes())
		require.Len(t, frames, 1)
		assert.Equal(t, "validation", frames[0].Event)

		problems, ok := frames[0].Data["consistencyProblems"].([]any)
		require.True(t, ok)
		require.Len(t, problems, 1)

		problem, ok := problems[0].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "domainNames", problem["fieldPath"])
	})

	t.Run("streams an error when the operation panics", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest("POST", "/", nil)

		Serve(ctx, func(context.Context) (any, error) {
			panic("boom")
		})

		frames := parseStreamFrames(t, recorder.Body.Bytes())
		require.Len(t, frames, 1)
		assert.Equal(t, "error", frames[0].Event)
		assert.Equal(t, "boom", frames[0].Data["message"])
	})

	t.Run("stops when the request is canceled", func(t *testing.T) {
		started := make(chan struct{})
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		requestContext, cancel := context.WithCancel(context.Background())
		ctx.Request = httptest.NewRequest("POST", "/", nil).WithContext(requestContext)

		done := make(chan struct{})
		go func() {
			defer close(done)

			Serve(ctx, func(runContext context.Context) (any, error) {
				close(started)
				<-runContext.Done()

				return nil, runContext.Err()
			})
		}()

		<-started
		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("Serve did not stop after the request context was canceled")
		}
	})
}

func parseStreamFrames(t *testing.T, body []byte) []streamFrame {
	t.Helper()

	frames := make([]streamFrame, 0)
	normalized := strings.ReplaceAll(string(body), "\r\n", "\n")

	for _, block := range strings.Split(normalized, "\n\n") {
		block = strings.TrimSpace(block)
		if block == "" || strings.HasPrefix(block, ":") {
			continue
		}

		frame := streamFrame{Data: map[string]any{}}
		dataLines := make([]string, 0)

		for _, line := range strings.Split(block, "\n") {
			if name, found := strings.CutPrefix(line, "event:"); found {
				frame.Event = strings.TrimSpace(name)
			}

			if data, found := strings.CutPrefix(line, "data:"); found {
				dataLines = append(dataLines, strings.TrimSpace(data))
			}
		}

		require.NotEmpty(t, dataLines)
		require.NoError(t, json.Unmarshal([]byte(strings.Join(dataLines, "\n")), &frame.Data))
		frames = append(frames, frame)
	}

	return frames
}
