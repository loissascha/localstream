package background

import (
	"context"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/media"
)

func (s *BackgroundService) createMovieStreams(ctx context.Context, movie *entity.Movie) error {
	probe, err := media.ProbeFile(movie.Path)
	if err != nil {
		return err
	}
	mFile, err := probe.ToMediaFile()
	if err != nil {
		return err
	}
	for _, v := range mFile.VideoStreams {
		vs := &entity.MovieVideoStream{
			MovieID:       movie.ID,
			Index:         v.Index,
			Profile:       v.Profile,
			Level:         v.Level,
			CodecName:     v.CodecName,
			CodecLongName: v.CodecLongName,
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
