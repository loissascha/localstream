package background

import (
	"time"

	"github.com/loissascha/localstream/internal/provider"
)

type seasonMetadataCache struct {
	created  time.Time
	metadata []provider.SeasonMetadata
}

func (s *BackgroundService) getSeasonMetadataResultLive(fetchID int) ([]provider.SeasonMetadata, error) {
	seasonMetadataResult, err := s.tvMetadataProvider.SearchSeasons(fetchID)
	s.seasonMetadataCache[fetchID] = seasonMetadataCache{
		created:  time.Now().UTC(),
		metadata: seasonMetadataResult,
	}
	return seasonMetadataResult, err
}

func (s *BackgroundService) getSeasonMetadataResultCacheOrLive(fetchID int) ([]provider.SeasonMetadata, error) {
	cachefile, ok := s.seasonMetadataCache[fetchID]
	if ok {
		if time.Now().UTC().Sub(cachefile.created) > 24*time.Hour {
			return s.getSeasonMetadataResultLive(fetchID)
		}
		return cachefile.metadata, nil
	}

	return s.getSeasonMetadataResultLive(fetchID)
}
