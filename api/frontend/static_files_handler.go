package frontend

import (
	"crypto/sha256"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"encoding/hex"
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

type staticFilesHandler struct {
	basePath *string
}

func (h staticFilesHandler) handle(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if h.basePath == nil || strings.HasPrefix(path, "/api/") {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	sanitizedPath, err := sanitizePath(path)
	if err != nil {
		log.Warnf("Request rejected. Possible path traversal attempt: %s", path)
		ctx.Status(http.StatusBadRequest)
		return
	}

	fileContents, fileType, err := h.loadFile(*sanitizedPath)
	if err != nil {
		log.Warnf("Error loading file %s: %s", *sanitizedPath, err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	etag, err := generateETag(fileType, fileContents)
	if err != nil {
		log.Warnf("Unable to generate ETag for file %s: %s", *sanitizedPath, err)
	} else if etag != "" {
		ifNoneMatchHeader := ctx.GetHeader("if-none-match")
		if ifNoneMatchHeader == etag {
			ctx.Status(http.StatusNotModified)
			return
		}

		ctx.Header("Cache-Control", "max-age=604800, must-revalidate")
		ctx.Header("ETag", etag)
	}

	ctx.Data(http.StatusOK, *fileType, fileContents)
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

func (h staticFilesHandler) loadFile(path string) ([]byte, *string, error) {
	baseDir := os.DirFS(*h.basePath)
	file, err := baseDir.Open(path)
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid) {
		path = indexFile
		file, err = baseDir.Open(path)
	}

	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

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

func generateETag(fileType *string, contents []byte) (string, error) {
	if fileType != nil && strings.HasPrefix(*fileType, "text/html") {
		return "", nil
	}

	hasher := sha256.New()
	if _, err := hasher.Write(contents); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	etag := hex.EncodeToString(hash)
	return `"` + etag + `"`, nil
}
