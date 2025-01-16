package frontend

import (
	"dillmann.com.br/nginx-ignition/core/common/log"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	indexFile        = "index.html"
	fallbackMimeType = "application/octet-stream"
)

type handler struct {
	basePath *string
}

func (h handler) handle(context *gin.Context) {
	path := context.Request.URL.Path
	if h.basePath == nil || strings.HasPrefix(path, "/api/") {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	sanitizedPath, err := sanitizePath(path)
	if err != nil {
		log.Warnf("Request rejected. Possible path traversal attempt: %s", path)
		context.Status(http.StatusBadRequest)
		return
	}

	fileContents, fileType, err := h.loadFile(*sanitizedPath)
	if err != nil {
		log.Warnf("Error loading file %s: %s", *sanitizedPath, err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Data(http.StatusOK, *fileType, fileContents)
}

func sanitizePath(path string) (*string, error) {
	parsedPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(parsedPath.Path)
	if strings.Contains(cleanPath, "..") {
		return nil, err
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(absPath, "/") {
		absPath = absPath[1:]
	}

	return &absPath, nil
}

func (h handler) loadFile(path string) ([]byte, *string, error) {
	baseDir := os.DirFS(*h.basePath)
	file, err := baseDir.Open(path)
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid) {
		path = indexFile
		file, err = baseDir.Open(path)
	}

	if err != nil {
		return nil, nil, err
	}

	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = fallbackMimeType
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	return fileContents, &mimeType, nil
}
