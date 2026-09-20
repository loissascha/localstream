package background

import (
	"time"

	"github.com/loissascha/localstream/internal/provider"
)

type seasonMetadataCache struct {
	created  time.Time
	metadata []provider.SeasonMetadata
}

type episodeMetadataCache struct {
	created  time.Time
	metadata []provider.EpisodeMetadata
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

func (s *BackgroundService) getEpisodeMetadataResultLive(fetchID int) ([]provider.EpisodeMetadata, error) {
	result, err := s.tvMetadataProvider.SearchEpisodes(fetchID)
	s.episodeMetadataCache[fetchID] = episodeMetadataCache{
		created:  time.Now().UTC(),
		metadata: result,
	}
	return result, err
}

func (s *BackgroundService) getEpisodeMetadataResultCacheOrLive(fetchID int) ([]provider.EpisodeMetadata, error) {
	cachefile, ok := s.episodeMetadataCache[fetchID]
	if ok {
		if time.Now().UTC().Sub(cachefile.created) > 24*time.Hour {
			return s.getEpisodeMetadataResultLive(fetchID)
		}
		// logger.Debug(nil, "Load episode metadata from cache")
		return cachefile.metadata, nil
	}

	return s.getEpisodeMetadataResultLive(fetchID)
}
