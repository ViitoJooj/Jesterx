package service

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
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
	Bucket   string `json:"bucket"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type StorageService struct{}

func NewStorageService() *StorageService {
	return &StorageService{}
}

func bucketForMime(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "images"
	case strings.HasPrefix(mimeType, "video/"):
		return "videos"
	default:
		return "documents"
	}
}

func allowedMime(mimeType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true,
		"image/webp": true, "image/svg+xml": true,
		"video/mp4": true, "video/webm": true, "video/ogg": true,
		"application/pdf":  true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	return allowed[mimeType]
}

func (s *StorageService) Upload(file io.Reader, filename string, size int64) (*UploadResult, error) {
	if config.SupabaseURL == "" || config.SupabaseServiceKey == "" {
		return nil, fmt.Errorf("supabase não configurado")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if !allowedMime(mimeType) {
		return nil, fmt.Errorf("tipo de arquivo não permitido: %s", mimeType)
	}

	if size > MaxUploadSize {
		return nil, fmt.Errorf("arquivo muito grande (máximo 50MB)")
	}

	bucket := bucketForMime(mimeType)
	uniqueName := uuid.New().String() + ext
	storagePath := time.Now().Format("2006/01/") + uniqueName

	data, err := io.ReadAll(io.LimitReader(file, MaxUploadSize+1))
	if err != nil {
		return nil, fmt.Errorf("leitura do arquivo: %w", err)
	}
	if int64(len(data)) > MaxUploadSize {
		return nil, fmt.Errorf("arquivo muito grande (máximo 50MB)")
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", config.SupabaseURL, bucket, storagePath)
	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("criar request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.SupabaseServiceKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "false")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase upload error %d: %s", resp.StatusCode, string(body))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", config.SupabaseURL, bucket, storagePath)

	return &UploadResult{
		URL:      publicURL,
		Path:     storagePath,
		Bucket:   bucket,
		MimeType: mimeType,
		Size:     int64(len(data)),
	}, nil
}
