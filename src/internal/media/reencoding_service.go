package media

import (
	"time"

	"github.com/google/uuid"
)

type ReencodeStatus struct {
	Id      uuid.UUID
	Started bool
}

type ReencodingService struct {
	reencodeMovieQueue map[uuid.UUID]ReencodeStatus
}

func NewReencodingService() *ReencodingService {
	return &ReencodingService{
		reencodeMovieQueue: map[uuid.UUID]ReencodeStatus{},
	}
}

func (s *ReencodingService) StartBackground() {
	// automatic reencodings
	go func() {
		time.Sleep(30 * time.Second)
		for {
			time.Sleep(60 * time.Second)
		}
	}()
}

func (s *ReencodingService) AddMovieToReencodeQueue(movieID uuid.UUID) ReencodeStatus {
	m, found := s.reencodeMovieQueue[movieID]
	if found {
		return m
	}
	m = ReencodeStatus{
		Id:      movieID,
		Started: false,
	}
	s.reencodeMovieQueue[movieID] = m
	return m
}
