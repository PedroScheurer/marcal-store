package storage

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	KindVideo = "video"
	KindImage = "image"
)

var (
	videoExtensions = map[string]bool{".mp4": true, ".mov": true, ".webm": true, ".m4v": true}
	imageExtensions = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
)

type MediaStorage struct {
	rootDir string
}

func NewMediaStorage(rootDir string) (*MediaStorage, error) {
	for _, sub := range []string{"videos", "images"} {
		if err := os.MkdirAll(filepath.Join(rootDir, sub), 0o755); err != nil {
			return nil, fmt.Errorf("create upload dir %s: %w", sub, err)
		}
	}
	return &MediaStorage{rootDir: rootDir}, nil
}

func (s *MediaStorage) RootDir() string {
	return s.rootDir
}

func (s *MediaStorage) Save(kind string, originalName string, contentType string, r io.Reader) (publicPath string, err error) {
	subdir, allowedExts, err := kindConfig(kind)
	if err != nil {
		return "", err
	}

	ext := pickExtension(originalName, contentType, allowedExts)
	fileName := uuid.NewString() + ext
	destPath := filepath.Join(s.rootDir, subdir, fileName)

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("write file: %w", err)
	}

	return fmt.Sprintf("/media/%s/%s", subdir, fileName), nil
}

func kindConfig(kind string) (subdir string, allowed map[string]bool, err error) {
	switch kind {
	case KindVideo:
		return "videos", videoExtensions, nil
	case KindImage:
		return "images", imageExtensions, nil
	default:
		return "", nil, fmt.Errorf("tipo inválido: use video ou image")
	}
}

func pickExtension(originalName, contentType string, allowed map[string]bool) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	if allowed[ext] {
		return ext
	}

	if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
		for _, candidate := range exts {
			candidate = strings.ToLower(candidate)
			if allowed[candidate] {
				return candidate
			}
		}
	}

	switch {
	case allowed[".mp4"]:
		return ".mp4"
	default:
		return ".jpg"
	}
}
