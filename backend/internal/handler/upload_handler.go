package handler

import (
	"io"
	"net/http"
	"strings"
	"log/slog"

	"wander/backend/internal/middleware"
	"wander/backend/internal/storage"
	"wander/backend/internal/utils"
)

const maxUploadBytes = 10 << 20 // 10 MB

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// UploadHandler accepts authenticated image uploads and hands the bytes to
// a storage.Provider (local disk or S3-compatible). It does not know or
// care where the file ends up — the returned URL comes straight from the
// provider, so URLs are absolute bucket URLs in S3 mode and server-relative
// ("/uploads/...") in local mode.
type UploadHandler struct {
	provider storage.Provider
}

// NewUploadHandler wires the handler to a storage.Provider.
func NewUploadHandler(provider storage.Provider) *UploadHandler {
	return &UploadHandler{provider: provider}
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

	// Sniff the first 512 bytes to detect the true MIME type. We read into
	// a buffer and then prepend it to the stream we pass to the provider so
	// the provider sees the full file, not from byte 513 onward.
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		utils.SendError(w, http.StatusBadRequest, "No se pudo leer el archivo", err.Error())
		return
	}
	mimeType := http.DetectContentType(buf[:n])

	ext, ok := allowedImageTypes[mimeType]
	if !ok {
		utils.SendError(w, http.StatusUnsupportedMediaType, "Tipo no soportado. Use JPG, PNG, WEBP o GIF", mimeType)
		return
	}

	// Reassemble: prepend the bytes we already consumed so the provider sees
	// a complete stream starting at byte 0.
	body := io.MultiReader(strings.NewReader(string(buf[:n])), file)

	result, err := h.provider.Save(r.Context(), body, mimeType, ext)
	if err != nil {
		slog.Error("upload failed", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Error al guardar el archivo", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusCreated, "Imagen subida", map[string]string{
		"url":      result.URL,
		"filename": result.Key,
		"original": sanitizeName(header.Filename),
	})
}

func sanitizeName(name string) string {
	if i := strings.LastIndex(name, "."); i > 0 {
		name = name[:i]
	}
	return name
}
