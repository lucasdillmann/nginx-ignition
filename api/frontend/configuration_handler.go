package frontend

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type configurationHandler struct {
	configuration *configuration.Configuration
}

func (h *configurationHandler) handle(ctx *gin.Context) {
	codeEditorApiKey, _ := h.configuration.Get("nginx-ignition.frontend.code-editor-api-key")
	version, _ := h.configuration.Get("nginx-ignition.version")

	var apiKey *string
	if codeEditorApiKey != "" {
		apiKey = &codeEditorApiKey
	}

	var versionString *string
	if version != "" {
		versionString = &version
	}

	output := &configurationDto{
		Version: versionDto{
			Current: versionString,
			Latest:  resolveLatestAvailableVersion(),
		},
		CodeEditor: codeEditorDto{
			ApiKey: apiKey,
		},
	}

	ctx.JSON(http.StatusOK, output)
}

func resolveLatestAvailableVersion() *string {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get("https://api.github.com/repos/lucasdillmann/nginx-ignition/releases?per_page=1&page=0")
	if err != nil {
		log.Warnf("Failed to fetch latest available version: %s", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read latest available version: %s", err)
		return nil
	}

	var releases []map[string]interface{}
	if err := json.Unmarshal(body, &releases); err != nil {
		log.Warnf("Failed to parse latest available version: %s", err)
		return nil
	}

	if len(releases) > 0 {
		if name, ok := releases[0]["name"].(string); ok {
			return &name
		}
	}

	return nil
}
