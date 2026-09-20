package background

import (
	"context"
	"fmt"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
)

func (s *BackgroundService) runSeasonsMatcher() error {
	// get all the seasons from the repo that have a show with a valid fetch source (not none/empty/multiple)
	return nil
}

func (s *BackgroundService) fetchMetadataForSeason(ctx context.Context, season *entity.Season, show *entity.Show) error {
	if !season.FetchSource.IsNone() {
		return nil
	}
	var err error
	if show == nil {
		show, err = s.showRepo.GetByID(ctx, season.ShowID)
		if err != nil {
			return err
		}
	}

	// show has invalid fetch source -> skip
	if !show.FetchSource.IsValid() {
		return nil
	}

	showMetadata, err := s.showMetadataRepo.GetByShowID(ctx, show.ID)
	if err != nil {
		return err
	}
	if len(showMetadata) != 1 {
		return fmt.Errorf("show has wrong amount of metadatas: %d", len(showMetadata))
	}

	seasonMetadataResult, err := s.getSeasonMetadataResultCacheOrLive(showMetadata[0].FetchID)
	if err != nil {
		return err
	}

	matchCount := 0
	for _, smr := range seasonMetadataResult {
		if smr.Number == season.Number {
			matchCount++
			err := s.seasonMetadataService.CreateSeasonMetadata(ctx, season, &smr)
			if err != nil {
				logger.Error(err, "[SeasonMatcher] Error creating metadata for season")
			}
			break
		}
	}

	switch matchCount {
	case 0:
		season.FetchSource = entity.FetchSourceEmpty
	case 1:
		season.FetchSource = entity.FetchSourceTVMaze
	case 2:
		season.FetchSource = entity.FetchSourceMultiple
	}

	err = s.seasonRepo.UpdateFetchSource(ctx, season.ID, season.FetchSource)
	if err != nil {
		return err
	}
	return nil
}
