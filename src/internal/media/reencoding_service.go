package media

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
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
			for id, st := range s.reencodeMovieQueue {
				if st.Started {
					continue
				}
				err := s.startMovieReencode(id)
				if err != nil {
					logger.Error(err, "Error running movie reencoding")
				}
				break
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

func (s *ReencodingService) startMovieReencode(id uuid.UUID) error {
	st, found := s.reencodeMovieQueue[id]
	if !found {
		return fmt.Errorf("movie not found")
	}
	st.Started = true
	s.reencodeMovieQueue[id] = st

	return nil
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
