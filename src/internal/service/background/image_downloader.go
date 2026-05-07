package background

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func downloadImageAndGetStaticPath(url string, filename string) (string, error) {
	ext, err := getExtensionFromUrl(url)
	if err != nil {
		return "", err
	}
	filename = filename + ext
	err = downloadImage(url, filename)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/static/%s", filename), nil
}

func downloadImage(url string, filename string) error {
	err := os.MkdirAll("./static", os.ModePerm)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %s", resp.Status)
	}

	filePath := filepath.Join("./static", filename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func getExtensionFromUrl(rawUrl string) (string, error) {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	ext := path.Ext(parsed.Path)
	if ext == "" {
		return "", fmt.Errorf("url has no file extension")
	}

	return ext, nil
}
