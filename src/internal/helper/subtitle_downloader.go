package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadSubtitleForEpisode(rawUrl string, showIdStr string, seasonNumberStr string, episodeNumberStr string) (string, error) {
	dirStr := "./static/subtitles/shows/" + showIdStr + "/" + seasonNumberStr + "/" + episodeNumberStr
	err := os.MkdirAll(dirStr, os.ModePerm)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(rawUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status %s", resp.Status)
	}

	filename, err := FilenameFromUrl(rawUrl)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dirStr, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func DownloadSubtitleForMovie(rawUrl string, movieIdStr string) (string, error) {
	dirStr := "./static/subtitles/movies/" + movieIdStr
	err := os.MkdirAll(dirStr, os.ModePerm)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(rawUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status %s", resp.Status)
	}

	filename, err := FilenameFromUrl(rawUrl)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(dirStr, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
