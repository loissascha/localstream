package helper

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertToVTT(inputPath string) (string, error) {
	if strings.TrimSpace(inputPath) == "" {
		return "", fmt.Errorf("input path is empty")
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext != ".srt" {
		return "", fmt.Errorf("expected .srt file, got %q", ext)
	}

	dir := filepath.Dir(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	outputPath := filepath.Join(dir, base+".vtt")

	cmd := exec.Command(
		"ffmpeg",
		"-y", // overwrite output if it exists
		"-i", inputPath,
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg convert %q to %q failed: %w\n%s", inputPath, outputPath, err, string(output))
	}

	return outputPath, nil
}
