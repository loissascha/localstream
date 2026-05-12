package helper

import (
	"fmt"
	"net/url"
	"path"
)

func GetExtensionFromUrl(rawUrl string) (string, error) {
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

func FilenameFromUrl(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	filename := path.Base(u.Path)
	if filename == "." || filename == "/" {
		return "", fmt.Errorf("url has no filename")
	}

	return filename, nil
}
