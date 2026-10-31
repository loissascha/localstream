package media

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type IntegrityCheckService struct {
	libRepo repository.LibraryRepository

	showRepo    repository.ShowRepository
	seasonRepo  repository.SeasonRepository
	movieRepo   repository.MovieRepository
	episodeRepo repository.EpisodeRepository

	movieIntegrityCheckIssueRepo repository.MovieIntegrityCheckIssueRepository
}

func NewIntegrityCheckService(
	libRepo repository.LibraryRepository,

	showRepo repository.ShowRepository,
	seasonRepo repository.SeasonRepository,
	movieRepo repository.MovieRepository,
	episodeRepo repository.EpisodeRepository,

	movieIntegrityCheckIssueRepo repository.MovieIntegrityCheckIssueRepository,
) *IntegrityCheckService {
	return &IntegrityCheckService{
		libRepo: libRepo,

		showRepo:    showRepo,
		seasonRepo:  seasonRepo,
		movieRepo:   movieRepo,
		episodeRepo: episodeRepo,

		movieIntegrityCheckIssueRepo: movieIntegrityCheckIssueRepo,
	}
}

func (s *IntegrityCheckService) StartBackground() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			logger.Info(nil, "Starting Integrity Check")
			err := s.RunOnce(ctx)
			if err != nil {
				slog.Error("error appeared while running integrity check service", "err", err)
			}
			logger.Info(nil, "Finished Integrity Check")
			time.Sleep(120 * time.Second)
		}
	}()
}

func (s *IntegrityCheckService) RunOnce(ctx context.Context) error {

	// search for an episode or movie that hasn't had the integrity check yet
	movie, err := s.movieRepo.FirstWithoutIntegrityChecked(ctx)
	if err != nil {
		return err
	}

	// do integrity check on movie
	if movie != nil {
		err := s.runOnMovie(ctx, movie)
		if err != nil {
			return err
		}
	}

	// do the integrity check on that file

	// record findings

	return nil
}

func (s *IntegrityCheckService) runOnMovie(ctx context.Context, movie *entity.Movie) error {
	result, err := Scan(ctx, movie.Path)
	if err != nil {
		return err
	}

	// log error and do nothing
	if result.ExitCode == -1 {
		slog.ErrorContext(ctx, "We were not able to scan the movie file for integrity check", "movie_id", movie.ID.String(), "movie_path", movie.Path)
		return nil
	}

	// log errors
	if result.Errors > 0 {
		slog.WarnContext(ctx, "Scanned movie has resulted in errors.", "movie_id", movie.ID.String(), "movie_path", movie.Path, "error_count", result.Errors)
		uid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		err = s.movieIntegrityCheckIssueRepo.Create(ctx, &entity.MovieIntegrityCheckIssue{
			ID:         uid,
			MovieID:    movie.ID,
			IssueLevel: entity.IntegrityCheckIssueLevelError,
			Issue:      fmt.Sprintf("Scan has shown %d errors on the movie file at %s", result.Errors, movie.Path),
		})
		if err != nil {
			return err
		}
	}

	// log warnings
	if result.Warnings > 0 {
		slog.WarnContext(ctx, "Scanned movie has resulted in warnings.", "movie_id", movie.ID.String(), "movie_path", movie.Path, "error_count", result.Warnings)
		uid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		err = s.movieIntegrityCheckIssueRepo.Create(ctx, &entity.MovieIntegrityCheckIssue{
			ID:         uid,
			MovieID:    movie.ID,
			IssueLevel: entity.IntegrityCheckIssueLevelWarning,
			Issue:      fmt.Sprintf("Scan has shown %d warnings on the movie file at %s", result.Errors, movie.Path),
		})
		if err != nil {
			return err
		}
	}

	// update integrity checked on the movie
	err = s.movieRepo.UpdateIntegrityChecked(ctx, movie.ID, true)
	if err != nil {
		return err
	}

	return nil
}
