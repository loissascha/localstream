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

func (self *MovieMatcher) RunBackground() {
	go func() {
		for {
			movie := <-self.Channel
			fmt.Println("movie", movie)
		}
	}()
}
