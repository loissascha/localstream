package background

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadImageAndGetStaticPath(originalUrl string, filename string) (string, error) {
	err := downloadImage(originalUrl, filename)
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
