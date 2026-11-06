package media

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/repository"
)

type ReencodeStatus struct {
	Id       uuid.UUID
	Started  bool
	Finished bool
	Failed   bool
	Error    string
}

type ReencodingService struct {
	movieRepo          repository.MovieRepository
	movieIntegrityRepo repository.MovieIntegrityCheckIssueRepository

	baseDir string

	movieQueueMutex    sync.Mutex
	reencodeMovieQueue map[uuid.UUID]ReencodeStatus
}

func NewReencodingService(
	baseDir string,
	movieRepo repository.MovieRepository,
	movieIntegrityRepo repository.MovieIntegrityCheckIssueRepository,
) *ReencodingService {
	return &ReencodingService{
		baseDir:            baseDir,
		movieRepo:          movieRepo,
		movieIntegrityRepo: movieIntegrityRepo,

		reencodeMovieQueue: map[uuid.UUID]ReencodeStatus{},
	}
}

func (s *ReencodingService) StartBackground() {
	// automatic reencodings
	go func() {
		time.Sleep(30 * time.Second)
		for {
			logger.Info(nil, "Starting Reencoding run")
			ctx, cancel := context.WithTimeout(context.Background(), 240*time.Minute)
			for id, st := range s.reencodeMovieQueue {
				if st.Started {
					continue
				}
				st.Started = true
				err := s.startMovieReencode(ctx, id)
				if err != nil {
					logger.Error(err, "Error running movie reencoding")
					st.Failed = true
					st.Error = err.Error()
					s.reencodeMovieQueue[id] = st
					continue
				}
				st.Finished = true
				s.reencodeMovieQueue[id] = st
				break
			}
			logger.Info(nil, "Finished Reencoding run")
			cancel()
			time.Sleep(60 * time.Second)
		}
	}()
}

func (s *ReencodingService) startMovieReencode(ctx context.Context, id uuid.UUID) error {
	st, found := s.reencodeMovieQueue[id]
	if !found {
		return fmt.Errorf("movie not found")
	}
	st.Started = true
	s.reencodeMovieQueue[id] = st

	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	p := filepath.Join(s.baseDir, fmt.Sprintf("%s.mp4", movie.ID.String()))
	err = ReencodeFile(ctx, movie.Path, p)
	if err != nil {
		return err
	}

	err = os.Rename(p, movie.Path)
	if err != nil {
		return err
	}

	err = s.movieIntegrityRepo.DeleteByMovieID(ctx, id)
	if err != nil {
		return err
	}

	err = s.movieRepo.UpdateIntegrityChecked(ctx, id, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReencodingService) AddMovieToReencodeQueue(movieID uuid.UUID) ReencodeStatus {
	s.movieQueueMutex.Lock()
	defer s.movieQueueMutex.Unlock()

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

func (s *ReencodingService) MovieQueue() map[uuid.UUID]ReencodeStatus {
	s.movieQueueMutex.Lock()
	defer s.movieQueueMutex.Unlock()

	snapshot := make(map[uuid.UUID]ReencodeStatus, len(s.reencodeMovieQueue))
	for id, st := range s.reencodeMovieQueue {
		snapshot[id] = st
	}

	return snapshot
}
