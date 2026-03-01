package i18n

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	corei18n "dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/i18n"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getDictionaryHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns bad request if language is invalid", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := corei18n.NewMockedCommands(ctrl)
			commands.EXPECT().GetDictionaries().Return([]i18n.Dictionary{})

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)

			handler := getDictionaryHandler{commands: commands}
			engine.GET("/:language", handler.handle)
			req, _ := http.NewRequest(http.MethodGet, "/invalid-language-tag", nil)
			engine.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("returns not found if dictionary is not available", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := corei18n.NewMockedCommands(ctrl)
			commands.EXPECT().GetDictionaries().Return([]i18n.Dictionary{})

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)

			handler := getDictionaryHandler{commands: commands}
			engine.GET("/:language", handler.handle)
			req, _ := http.NewRequest(http.MethodGet, "/en-US", nil)
			engine.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns dictionary as DTO if found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := corei18n.NewMockedCommands(ctrl)
			dict := newDictionary()
			dicts := []i18n.Dictionary{dict}
			commands.EXPECT().GetDictionaries().Return(dicts)

			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)

			handler := getDictionaryHandler{commands: commands}
			engine.GET("/:language", handler.handle)
			req, _ := http.NewRequest(http.MethodGet, "/"+dict.Language().String(), nil)
			engine.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response dictionaryResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, toDictionaryDTO(dict), response)
		})
	})
}
