package frontend

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

const (
	indexFile        = "index.html"
	fallbackMimeType = "application/octet-stream"
)

type staticFilesHandler struct {
	basePath *string
}

func (h staticFilesHandler) handle(ctx *gin.Context) {
	urlPath := ctx.Request.URL.Path
	if h.basePath == nil || strings.HasPrefix(urlPath, "/api/") {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	sanitizedPath, err := sanitizePath(urlPath)
	if err != nil {
		log.Warnf("Request rejected. Possible path traversal attempt: %s", urlPath)
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

func sanitizePath(p string) (*string, error) {
	if p == "" {
		return nil, errors.New("path cannot be empty")
	}

	cleanPath := path.Clean(p)
	if strings.Contains(cleanPath, "..") || strings.Contains(p, "..") {
		return nil, errors.New("invalid path")
	}

	cleanPath = strings.TrimPrefix(cleanPath, "/")
	if cleanPath == "" {
		cleanPath = indexFile
	}

	return &cleanPath, nil
}

func (h staticFilesHandler) loadFile(filePath string) ([]byte, *string, error) {
	baseDir := os.DirFS(*h.basePath)
	file, err := baseDir.Open(filePath)
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid) {
		filePath = indexFile
		file, err = baseDir.Open(filePath)
	}

	if err != nil {
		return nil, nil, err
	}

	//nolint:errcheck
	defer file.Close()

	ext := filepath.Ext(filePath)
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
