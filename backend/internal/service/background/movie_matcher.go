package background

import (
	"fmt"

	"github.com/loissascha/localstream/internal/entity"
)

type MovieMatcher struct {
	Channel chan *entity.Movie
}

func NewMovieMatcher() *MovieMatcher {
	ch := make(chan *entity.Movie)
	return &MovieMatcher{
		Channel: ch,
	}
}

func (l *MovieMatcher) RunBackground() {
	go func() {
		for {
			movie := <-l.Channel
			fmt.Println("movie", movie)
		}
	}()
}
