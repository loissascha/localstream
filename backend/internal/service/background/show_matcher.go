package background

import (
	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
)

type ShowMatcher struct {
	channel chan uuid.UUID
}

func NewShowMatcher() *ShowMatcher {
	ch := make(chan uuid.UUID)
	return &ShowMatcher{
		channel: ch,
	}
}

func (l *ShowMatcher) RunBackground() {
	go func() {
		for {
			showId := <-l.channel
			logger.Info(nil, "New ShowID triggered! {ShowId}", showId)
		}
	}()
}
