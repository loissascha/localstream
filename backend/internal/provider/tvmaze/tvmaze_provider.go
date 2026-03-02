package tvmaze

import "github.com/loissascha/go-logger/logger"

type TVMazeProvider struct {
}

func NewTVMazeProvider() *TVMazeProvider {
	return &TVMazeProvider{}
}

func (p *TVMazeProvider) SearchSeries(name string) {
	logger.Info(nil, "Search series {Name}", name)
}
