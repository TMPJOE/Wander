package handler

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"wander/backend/internal/middleware"
	"wander/backend/internal/utils"
)

const (
	maxUploadBytes = 10 << 20 // 10 MB
	uploadsSubdir  = "uploads"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadHandler struct {
	uploadsDir string
}

func NewUploadHandler(uploadsDir string) *UploadHandler {
	return &UploadHandler{uploadsDir: uploadsDir}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserID(r.Context()); !ok {
		utils.SendError(w, http.StatusUnauthorized, "No autorizado", nil)
		return
	}

	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Archivo demasiado grande o inválido (máx 10MB)", err.Error())
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Campo 'image' requerido", err.Error())
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	if _, err := file.ReadAt(buf, 0); err != nil {
		utils.SendError(w, http.StatusBadRequest, "No se pudo leer el archivo", err.Error())
		return
	}
	mimeType := http.DetectContentType(buf)

	ext, ok := allowedImageTypes[mimeType]
	if !ok {
		utils.SendError(w, http.StatusUnsupportedMediaType, "Tipo no soportado. Use JPG, PNG, WEBP o GIF", mimeType)
		return
	}

	if err := os.MkdirAll(h.uploadsDir, 0o755); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al crear directorio", err.Error())
		return
	}

	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al generar nombre", err.Error())
		return
	}
	filename := hex.EncodeToString(randBytes) + ext
	dstPath := filepath.Join(h.uploadsDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al guardar archivo", err.Error())
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al leer archivo", err.Error())
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error al escribir archivo", err.Error())
		return
	}

	url := "/uploads/" + filename
	utils.SendSuccess(w, http.StatusCreated, "Imagen subida", map[string]string{
		"url":      url,
		"filename": filename,
		"original": sanitizeName(header.Filename),
	})
}

func sanitizeName(name string) string {
	if i := strings.LastIndex(name, "."); i > 0 {
		name = name[:i]
	}
	return name
}
