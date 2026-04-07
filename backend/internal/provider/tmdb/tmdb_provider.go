package tmdb

import (
	"github.com/loissascha/localstream/internal/provider"
)

type TMDBProvider struct {
}

func NewTMDBProvider() *TMDBProvider {
	return &TMDBProvider{}
}

func (self *TMDBProvider) SearchMovie(name string, year int) {

}

var _ provider.MovieMetadataProvider = (*TMDBProvider)(nil)
