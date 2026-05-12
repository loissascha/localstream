package helper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipSingleFile(zipPath string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()

	outputDir := filepath.Dir(zipPath)

	var zipFile *zip.File

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if zipFile != nil {
			return "", fmt.Errorf("zip contains more than one file")
		}

		zipFile = file
	}

	if zipFile == nil {
		return "", fmt.Errorf("zip does not contain a file")
	}

	targetPath := filepath.Join(outputDir, zipFile.Name)

	// Prevent ZipSlip paths like ../../evil.txt
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return "", fmt.Errorf("get output dir abs path: %w", err)
	}

	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("get target abs path: %w", err)
	}

	if !strings.HasPrefix(absTargetPath, absOutputDir+string(os.PathSeparator)) {
		return "", fmt.Errorf("illegal file path in zip: %s", zipFile.Name)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return "", fmt.Errorf("create parent dir: %w", err)
	}

	src, err := zipFile.Open()
	if err != nil {
		return "", fmt.Errorf("open zipped file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("create target file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}

	return targetPath, nil
}
