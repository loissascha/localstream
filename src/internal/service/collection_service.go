package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidCollectionInput = errors.New("invalid collection input")

type CreateCollectionInput struct {
	Name string
}

type UpdateCollectionInput struct {
	Name string
}

type CollectionService struct {
	collectionRepo repository.CollectionRepository
}

func NewCollectionService(collectionRepo repository.CollectionRepository) *CollectionService {
	return &CollectionService{collectionRepo: collectionRepo}
}

func (s *CollectionService) Create(ctx context.Context, userID int64, input CreateCollectionInput) (*entity.Collection, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	name, err := normalizeCollectionName(input.Name)
	if err != nil {
		return nil, err
	}

	collection := &entity.Collection{
		UserID: userID,
		Name:   name,
	}

	if err := s.collectionRepo.Create(ctx, collection); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	return collection, nil
}

func (s *CollectionService) GetByID(ctx context.Context, userID int64, collectionID string) (*entity.Collection, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	collection, err := s.collectionRepo.GetByIDForUser(ctx, collectionUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get collection by id: %w", err)
	}

	return collection, nil
}

func (s *CollectionService) ListByUserID(ctx context.Context, userID int64) ([]entity.Collection, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	collections, err := s.collectionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list collections by user id: %w", err)
	}

	return collections, nil
}

func (s *CollectionService) UpdateName(ctx context.Context, userID int64, collectionID string, input UpdateCollectionInput) (*entity.Collection, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	name, err := normalizeCollectionName(input.Name)
	if err != nil {
		return nil, err
	}

	if err := s.collectionRepo.UpdateName(ctx, collectionUUID, userID, name); err != nil {
		return nil, fmt.Errorf("update collection name: %w", err)
	}

	collection, err := s.collectionRepo.GetByIDForUser(ctx, collectionUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("get collection after update: %w", err)
	}

	return collection, nil
}

func (s *CollectionService) DeleteByID(ctx context.Context, userID int64, collectionID string) error {
	if userID <= 0 {
		return ErrInvalidCollectionInput
	}

	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	if err := s.collectionRepo.DeleteByIDForUser(ctx, collectionUUID, userID); err != nil {
		return fmt.Errorf("delete collection by id: %w", err)
	}

	return nil
}

func (s *CollectionService) AddMovie(ctx context.Context, userID int64, collectionID, movieID string) error {
	if userID <= 0 {
		return ErrInvalidCollectionInput
	}

	collectionUUID, movieUUID, err := decodeCollectionAndMovieID(collectionID, movieID)
	if err != nil {
		return err
	}

	if err := s.collectionRepo.AddMovie(ctx, userID, collectionUUID, movieUUID); err != nil {
		return fmt.Errorf("add movie to collection: %w", err)
	}

	return nil
}

func (s *CollectionService) RemoveMovie(ctx context.Context, userID int64, collectionID, movieID string) error {
	if userID <= 0 {
		return ErrInvalidCollectionInput
	}

	collectionUUID, movieUUID, err := decodeCollectionAndMovieID(collectionID, movieID)
	if err != nil {
		return err
	}

	if err := s.collectionRepo.RemoveMovie(ctx, userID, collectionUUID, movieUUID); err != nil {
		return fmt.Errorf("remove movie from collection: %w", err)
	}

	return nil
}

func (s *CollectionService) AddShow(ctx context.Context, userID int64, collectionID, showID string) error {
	if userID <= 0 {
		return ErrInvalidCollectionInput
	}

	collectionUUID, showUUID, err := decodeCollectionAndShowID(collectionID, showID)
	if err != nil {
		return err
	}

	if err := s.collectionRepo.AddShow(ctx, userID, collectionUUID, showUUID); err != nil {
		return fmt.Errorf("add show to collection: %w", err)
	}

	return nil
}

func (s *CollectionService) RemoveShow(ctx context.Context, userID int64, collectionID, showID string) error {
	if userID <= 0 {
		return ErrInvalidCollectionInput
	}

	collectionUUID, showUUID, err := decodeCollectionAndShowID(collectionID, showID)
	if err != nil {
		return err
	}

	if err := s.collectionRepo.RemoveShow(ctx, userID, collectionUUID, showUUID); err != nil {
		return fmt.Errorf("remove show from collection: %w", err)
	}

	return nil
}

func (s *CollectionService) ListMovies(ctx context.Context, userID int64, collectionID string) ([]repository.MovieSelectItem, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	movies, err := s.collectionRepo.ListMovies(ctx, userID, collectionUUID)
	if err != nil {
		return nil, fmt.Errorf("list collection movies: %w", err)
	}

	return movies, nil
}

func (s *CollectionService) ListShows(ctx context.Context, userID int64, collectionID string) ([]repository.ShowSelectItem, error) {
	if userID <= 0 {
		return nil, ErrInvalidCollectionInput
	}

	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	shows, err := s.collectionRepo.ListShows(ctx, userID, collectionUUID)
	if err != nil {
		return nil, fmt.Errorf("list collection shows: %w", err)
	}

	return shows, nil
}

func normalizeCollectionName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", ErrInvalidCollectionInput
	}

	return trimmedName, nil
}

func decodeCollectionAndMovieID(collectionID, movieID string) (uuid.UUID, uuid.UUID, error) {
	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("decode movie id: %w", ErrInvalidCollectionInput)
	}

	return collectionUUID, movieUUID, nil
}

func decodeCollectionAndShowID(collectionID, showID string) (uuid.UUID, uuid.UUID, error) {
	collectionUUID, err := encoders.DecodeUUID(collectionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("decode collection id: %w", ErrInvalidCollectionInput)
	}

	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("decode show id: %w", ErrInvalidCollectionInput)
	}

	return collectionUUID, showUUID, nil
}
