package background

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/media"
)

func (s *BackgroundService) RunMediaStreamChecks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	movies, err := s.movieRepo.ListWithoutVideoStream(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	for _, m := range movies {
		err := s.createMovieStreams(ctx, &m)
		if err != nil {
			slog.Error("error creating movie stream", "err", err, "movieID", m.ID.String())
		}
	}

	// TODO: same for episodes

	return nil
}

func (s *BackgroundService) createMovieStreams(ctx context.Context, movie *entity.Movie) error {
	probe, err := media.ProbeFile(movie.Path)
	if err != nil {
		return err
	}
	mFile, err := probe.ToMediaFile()
	if err != nil {
		return err
	}

	splitStr := strings.Split(mFile.Filename, ".")
	container := splitStr[len(splitStr)-1]

	for _, v := range mFile.VideoStreams {
		vs := &entity.MovieVideoStream{
			MovieID:       movie.ID,
			Index:         v.Index,
			Container:     container,
			Profile:       v.Profile,
			Level:         v.Level,
			CodecName:     v.CodecName,
			CodecLongName: v.CodecLongName,
			PixelFormat:   v.PixelFormat,
			Duration:      v.Duration,
			Width:         v.Width,
			Height:        v.Height,
			Tags:          v.Tags,
			BitRate:       v.BitRate,
		}
		err := s.movieVideoStreamRepo.Create(ctx, vs)
		if err != nil {
			return err
		}
	}
	for _, a := range mFile.AudioStreams {
		as := &entity.MovieAudioStream{
			MovieID:       movie.ID,
			Index:         a.Index,
			Profile:       a.Profile,
			CodecName:     a.CodecName,
			CodecLongName: a.CodecLongName,
			Duration:      a.Duration,
			Tags:          a.Tags,
			BitRate:       a.BitRate,
		}
		err := s.movieAudioStreamRepo.Create(ctx, as)
		if err != nil {
			return err
		}
	}
	return nil
}
