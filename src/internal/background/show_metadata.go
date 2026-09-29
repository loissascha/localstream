package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/fetchsource"
)

func (s *BackgroundService) runShowsMatcher() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// get all the shows from the repo that have no metadata (fetchsource none)
	shows, err := s.showRepo.NecessaryForMetadataFetch(ctx)
	if err != nil {
		return err
	}
	slog.Info("shows matcher is testing how many shows", "amount", len(shows))

	for _, show := range shows {
		err := s.fetchMetadataForShow(ctx, &show)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BackgroundService) fetchMetadataForShow(ctx context.Context, show *entity.Show) error {
	if !show.FetchSource.IsNone() {
		return nil
	}

	showSearchResults, err := s.tvMetadataProvider.SearchShow(show.Name, show.Year)
	if err != nil {
		return err
	}

	for _, res := range showSearchResults {
		err := s.showMetadataService.CreateShowMetadata(ctx, show, &res.Show)
		if err != nil {
			return err
		}
	}

	if len(showSearchResults) > 1 {
		show.FetchSource = fetchsource.FetchSourceMultiple
		s.showRepo.UpdateFetchSource(ctx, show.ID, fetchsource.FetchSourceMultiple)
		logger.Info(nil, "Found multiple results for show {Show} ({Year})", show.Name, show.Year)
	} else if len(showSearchResults) == 1 {
		show.FetchSource = fetchsource.FetchSourceTVMaze
		s.showRepo.UpdateFetchSource(ctx, show.ID, fetchsource.FetchSourceTVMaze)
		logger.Info(nil, "Found perfect match for show {Show} ({Year}): {Match}", show.Name, show.Year, showSearchResults[0])
	} else {
		show.FetchSource = fetchsource.FetchSourceEmpty
		s.showRepo.UpdateFetchSource(ctx, show.ID, fetchsource.FetchSourceEmpty)
	}
	return nil
}
