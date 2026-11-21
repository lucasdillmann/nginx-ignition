package certificate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_issueHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("streams a successful result", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			certificateData := newCertificate()
			payload := newIssueCertificateRequest()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(certificateData, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates/issue",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "text/event-stream", recorder.Header().Get("Content-Type"))

			frame := lastFrame(t, parseSSE(t, recorder.Body.Bytes()))
			assert.Equal(t, "result", frame.Event)
			assert.Equal(t, true, frame.Data["success"])
			assert.Equal(t, certificateData.ID.String(), frame.Data["certificateId"])
		})

		t.Run("streams over HTTP", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			certificateData := newCertificate()
			payload := newIssueCertificateRequest()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(certificateData, nil)

			engine := gin.New()
			engine.POST(
				"/api/certificates/issue",
				extendWriteDeadline,
				issueHandler{commands}.handle,
			)

			server := httptest.NewServer(engine)
			defer server.Close()

			body, _ := json.Marshal(payload)
			response, err := http.Post(
				server.URL+"/api/certificates/issue",
				"application/json",
				bytes.NewBuffer(body),
			)
			require.NoError(t, err)
			defer response.Body.Close()

			require.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, "text/event-stream", response.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			frame := lastFrame(t, parseSSE(t, responseBody))
			assert.Equal(t, "result", frame.Event)
			assert.Equal(t, true, frame.Data["success"])
			assert.Equal(t, certificateData.ID.String(), frame.Data["certificateId"])
		})

		t.Run("streams progress events before the result", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			certificateData := newCertificate()
			payload := newIssueCertificateRequest()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				DoAndReturn(func(
					ctx context.Context,
					_ *certificate.IssueRequest,
				) (*certificate.Certificate, error) {
					progress.Report(ctx, progress.LevelDebug, "debug message")
					progress.Report(ctx, progress.LevelInfo, "info message")

					return certificateData, nil
				})

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates/issue",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			frames := parseSSE(t, recorder.Body.Bytes())
			require.Len(t, frames, 3)

			assert.Equal(t, "log", frames[0].Event)
			assert.Equal(t, "debug", frames[0].Data["level"])
			assert.Equal(t, "debug message", frames[0].Data["message"])

			assert.Equal(t, "log", frames[1].Event)
			assert.Equal(t, "info", frames[1].Data["level"])

			assert.Equal(t, "result", frames[2].Event)
			assert.Equal(t, true, frames[2].Data["success"])
		})

		t.Run("streams a failed result on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newIssueCertificateRequest()
			expectedErr := assert.AnError
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates/issue",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			frame := lastFrame(t, parseSSE(t, recorder.Body.Bytes()))
			assert.Equal(t, "result", frame.Event)
			assert.Equal(t, false, frame.Data["success"])
			assert.Equal(t, expectedErr.Error(), frame.Data["errorReason"])
		})

		t.Run("streams validation problems", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newIssueCertificateRequest()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(nil, validation.SingleFieldError(
					"domainNames",
					i18n.M(t.Context(), i18n.K.CommonCannotBeEmpty),
				))

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates/issue",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			frame := lastFrame(t, parseSSE(t, recorder.Body.Bytes()))
			assert.Equal(t, "validation", frame.Event)

			problems, ok := frame.Data["consistencyProblems"].([]any)
			require.True(t, ok)
			require.Len(t, problems, 1)

			problem, ok := problems[0].(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "domainNames", problem["fieldPath"])
		})

		t.Run("stops when the request is canceled", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			started := make(chan struct{})
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, _ *certificate.IssueRequest) (*certificate.Certificate, error) {
					close(started)
					<-ctx.Done()
					return nil, ctx.Err()
				})

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			requestContext, cancel := context.WithCancel(context.Background())
			payload := newIssueCertificateRequest()
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates/issue",
				bytes.NewBuffer(body),
			).WithContext(requestContext)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}

			done := make(chan struct{})
			go func() {
				defer close(done)
				handler.handle(ginContext)
			}()

			<-started
			cancel()

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("handler did not stop after the request context was canceled")
			}
		})
	})
}
