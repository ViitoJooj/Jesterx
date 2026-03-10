package service

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ViitoJooj/Jesterx/internal/config"
)

const MaxUploadSize = 50 * 1024 * 1024 // 50 MB

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type StorageService struct{}

func NewStorageService() *StorageService {
	return &StorageService{}
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	"video/mp4":     true,
	"video/webm":    true,
	"video/ogg":     true,
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// categoryForMime returns the sub-folder name based on the MIME type.
func categoryForMime(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "images"
	case strings.HasPrefix(mimeType, "video/"):
		return "videos"
	default:
		return "docs"
	}
}

func (s *StorageService) Upload(file io.Reader, filename string, size int64) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}

	if !allowedMimeTypes[mimeType] {
		return nil, fmt.Errorf("tipo de arquivo não permitido: %s", mimeType)
	}

	if size > MaxUploadSize {
		return nil, fmt.Errorf("arquivo muito grande (máximo 50MB)")
	}

	data, err := io.ReadAll(io.LimitReader(file, MaxUploadSize+1))
	if err != nil {
		return nil, fmt.Errorf("leitura do arquivo: %w", err)
	}
	if int64(len(data)) > MaxUploadSize {
		return nil, fmt.Errorf("arquivo muito grande (máximo 50MB)")
	}

	category := categoryForMime(mimeType)
	datePart := time.Now().UTC().Format("2006/01")
	uniqueName := uuid.New().String() + ext
	// Relative path inside the data folder: e.g. images/2025/01/<uuid>.jpg
	relPath := filepath.Join(category, filepath.FromSlash(datePart), uniqueName)
	absPath := filepath.Join(config.StoragePath, relPath)

	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return nil, fmt.Errorf("criar diretório de storage: %w", err)
	}

	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return nil, fmt.Errorf("salvar arquivo: %w", err)
	}

	// Public URL served by the backend at /files/...
	publicURL := "/files/" + filepath.ToSlash(relPath)

	return &UploadResult{
		URL:      publicURL,
		Path:     filepath.ToSlash(relPath),
		MimeType: mimeType,
		Size:     int64(len(data)),
	}, nil
}

