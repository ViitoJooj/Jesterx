package handlers

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "field 'file' is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := h.storageService.Upload(file, header.Filename, header.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "upload realizado",
		"data":    result,
	})
}
