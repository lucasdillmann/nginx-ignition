package i18n

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/i18n/dict"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getDictionaryHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns dictionaries as DTOs", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := i18n.NewMockedCommands(ctrl)
			dicts := []dict.Dictionary{newDictionary()}
			defaultLanguage := language.BrazilianPortuguese
			commands.EXPECT().GetDictionaries().Return(dicts)
			commands.EXPECT().DefaultLanguage().Return(defaultLanguage)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			handler := getDictionaryHandler{commands: commands}
			handler.handle(ctx)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response dictionariesDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, toDTO(defaultLanguage, dicts), response)
		})
	})
}
