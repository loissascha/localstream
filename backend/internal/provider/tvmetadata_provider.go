package provider

import "github.com/loissascha/localstream/internal/parsers"

type TVMetadataProvider interface {
	SearchSeries(episodeInfo *parsers.EpisodeInfo)
}
