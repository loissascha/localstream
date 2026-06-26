package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/encoders"
)

// TODO: update so that files go to their respective paths

func GetShowImagePath(showID uuid.UUID) string {
	return fmt.Sprintf("images/shows/%s", encoders.EncodeUUID(showID))
}

func GetMovieImagePath(movieID uuid.UUID) string {
	return fmt.Sprintf("images/movies/%s", encoders.EncodeUUID(movieID))
}

func DownloadImageAndGetStaticPath(url string, pathPrefix string, filename string) (string, error) {
	ext, err := GetExtensionFromUrl(url)
	if err != nil {
		return "", err
	}
	filename = filename + ext
	fp, err := downloadImage(url, pathPrefix, filename)
	if err != nil {
		return "", err
	}
	return fp, nil
}

func downloadImage(url string, pathPrefix string, filename string) (string, error) {
	basepath := filepath.Join("./static", pathPrefix)
	err := os.MkdirAll(basepath, os.ModePerm)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %s", resp.Status)
	}

	filePath := filepath.Join(basepath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(basepath, "/") {
		basepath = "/" + basepath
	}
	return filePath, nil
}
