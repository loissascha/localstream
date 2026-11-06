package media

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type ScanResult struct {
	Errors   int
	Warnings int
	Samples  []string // first 20 warning/error messages
	ExitCode int      // 0, 69, or another process exit code; -1 if not started
}

func Scan(ctx context.Context, path string) (ScanResult, error) {
	result := ScanResult{ExitCode: -1}

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner", "-nostdin", "-nostats",
		"-loglevel", "repeat+level+warning",
		"-i", path,
		"-map", "0:V?", "-map", "0:a?",
		"-f", "null", "-",
	)
	cmd.Stdout = io.Discard
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return result, err
	}

	if err := cmd.Start(); err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.Contains(line, "[error]"),
			strings.Contains(line, "[fatal]"),
			strings.Contains(line, "[panic]"):
			result.Errors++
		case strings.Contains(line, "[warning]"):
			result.Warnings++
		default:
			continue
		}

		if len(result.Samples) < 20 {
			result.Samples = append(result.Samples, line)
		}
	}

	readErr := scanner.Err()
	waitErr := cmd.Wait()

	if ctx.Err() != nil {
		return result, ctx.Err() // cancelled / timed out: NOT a clean scan
	}
	if readErr != nil {
		return result, fmt.Errorf("read ffmpeg log: %w", readErr)
	}
	if waitErr != nil {
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			result.ExitCode = exitErr.ExitCode()

			// FFmpeg reserves 69 for excessive decode failures;
			// it normally processes the input to completion first.
			if result.ExitCode == 69 {
				return result, nil
			}
		}
		return result, fmt.Errorf("ffmpeg scan failed: %w", waitErr)
	}

	result.ExitCode = 0
	return result, nil
}

func ReencodeFile(ctx context.Context, path string, output string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner", "-nostdin", "-n",
		"-i", path,
		"-map", "0:V:0", "-map", "0:a?",
		"-c:v", "libx264", "-preset", "medium", "-crf", "20",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac", "-b:a", "192k",
		"-movflags", "+faststart", output,
	)
	w := &strings.Builder{}
	cmd.Stdout = io.Discard
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg reencode failed: %w: %s", err, w.String())
	}
	return nil
}
