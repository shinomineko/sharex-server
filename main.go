package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	uploadDir            = "./uploads"
	defaultMaxUploadSize = 20
)

var (
	uploadKey     string
	maxUploadSize int64
)

func main() {
	uploadKey = os.Getenv("SHAREX_UPLOAD_KEY")
	if uploadKey == "" {
		log.Fatal("SHAREX_UPLOAD_KEY must be set")
	}

	sizeStr := os.Getenv("SHAREX_MAX_UPLOAD_SIZE_MB")
	if sizeStr == "" {
		maxUploadSize = defaultMaxUploadSize << 20
	} else {
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			log.Printf("Invalid SHAREX_MAX_UPLOAD_SIZE_MB value '%s', using default 20MB", sizeStr)
			maxUploadSize = defaultMaxUploadSize << 20
		} else {
			maxUploadSize = size << 20
		}
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", authMiddleware(handleUpload))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(uploadDir))))

	log.Printf("Server starting on port 3939")
	if err := http.ListenAndServe(":3939", nil); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != uploadKey {
			http.Error(w, "Invalid upload key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func generateFilename(originalName string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	ext := filepath.Ext(originalName)
	return hex.EncodeToString(bytes) + ext, nil
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := generateFilename(header.Filename)
	if err != nil {
		http.Error(w, "Error generating filename", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("%s://%s/files/%s",
		getProtocol(r),
		r.Host,
		filename,
	)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"success": true,
		"file": {
			"url": "%s",
			"name": "%s",
			"size": %d
		}
	}`, fileURL, filename, header.Size)
}

func getProtocol(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
