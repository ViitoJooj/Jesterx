package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type StorageHandler struct {
	storageService *service.StorageService
}

func NewStorageHandler(s *service.StorageService) *StorageHandler {
	return &StorageHandler{storageService: s}
}

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 50MB limit (matches service constant)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		jsonError(w, "falha ao processar formulário: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonError(w, "campo 'file' é obrigatório", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := h.storageService.Upload(file, header.Filename, header.Size)
	if err != nil {
		status := http.StatusBadRequest
		msg := err.Error()
		if strings.Contains(msg, "autenticação") || strings.Contains(msg, "SERVICE_KEY") {
			status = http.StatusInternalServerError
		}
		jsonError(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "upload realizado com sucesso",
		"data":    result,
	})
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
}
