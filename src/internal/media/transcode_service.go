package media

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type TranscodeService struct {
	baseDir string
}

func NewTranscodeService(baseDir string) *TranscodeService {
	return &TranscodeService{
		baseDir: baseDir,
	}
}

func (s *TranscodeService) StartTranscode(path string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	slog.Info("starting new transcode", "id", id)

	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	return nil
}
