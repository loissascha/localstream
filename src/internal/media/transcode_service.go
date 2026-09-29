package media

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

type TranscodeSession struct {
	ID        string
	InputPath string
	OutputDir string
	Cmd       *exec.Cmd
}

type TranscodeService struct {
	baseDir         string
	runningSessions map[string]*TranscodeSession
}

func NewTranscodeService(baseDir string) *TranscodeService {
	return &TranscodeService{
		baseDir:         baseDir,
		runningSessions: map[string]*TranscodeSession{},
	}
}

func (s *TranscodeService) StopTranscode(id string) error {
	sess, found := s.runningSessions[id]
	if !found {
		return errors.New("session not found")
	}

	err := sess.Cmd.Cancel()
	if err != nil {
		return err
	}

	err = os.Remove(sess.OutputDir)
	if err != nil {
		return err
	}

	delete(s.runningSessions, id)

	return nil
}

func (s *TranscodeService) StartTranscode(path string) (*TranscodeSession, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	idstr := id.String()
	slog.Info("starting new transcode", "id", idstr)

	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}

	outputDir := filepath.Join(s.baseDir, idstr)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"ffmpeg",

		"-i", path,

		"-map", "0:v:0",
		"-map", "0:a:0?",

		"-c:v", "libx264",
		"-preset", "veryfast",
		"-profile:v", "high",
		"-level:v", "4.0",
		"-pix_fmt", "yuv420p",
		"-crf", "23",

		"-c:a", "aac",
		"-b:a", "192k",
		"-ac", "2",

		"-f", "hls",
		"-hls_time", "4",
		"-hls_list_size", "0",
		"-hls_segment_type", "fmp4",

		"-hls_fmp4_init_filename", "init.mp4",
		"-hls_segment_filename",
		filepath.Join(outputDir, "segment_%05d.m4s"),
		filepath.Join(outputDir, "stream.m3u8"),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sess := &TranscodeSession{
		ID:        idstr,
		InputPath: path,
		OutputDir: outputDir,
		Cmd:       cmd,
	}
	s.runningSessions[idstr] = sess

	return sess, nil
}
