package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/repository"
	"github.com/loissascha/localstream/internal/service"
)

type CollectionHandler struct {
	s                 *server.Server
	authMiddleware    *middleware.AuthMiddleware
	collectionService *service.CollectionService
}

type CreateCollectionRequest struct {
	Name string `json:"name"`
}

type UpdateCollectionRequest struct {
	Name string `json:"name"`
}

func NewCollectionHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, collectionService *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		s:                 s,
		authMiddleware:    authMiddleware,
		collectionService: collectionService,
	}
}

func (h *CollectionHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/collections",
		h.listCollections,
		server.WithExportType[CollectionListResponse](),
		server.WithExportType[CollectionInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/collections",
		h.createCollection,
		server.WithExportType[CreateCollectionRequest](),
		server.WithExportType[CollectionInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GETI("/api/v1/collections/{collectionID}",
		h.getCollection,
		server.WithExportType[CollectionDetailResponse](),
		server.WithExportType[CollectionInfo](),
		server.WithExportType[MovieInfo](),
		server.WithExportType[ShowInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/collections/{collectionID}/update",
		h.updateCollection,
		server.WithExportType[UpdateCollectionRequest](),
		server.WithExportType[CollectionInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.DELETEI("/api/v1/collections/{collectionID}",
		h.deleteCollection,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/collections/{collectionID}/movies/{movieID}",
		h.addMovie,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.DELETEI("/api/v1/collections/{collectionID}/movies/{movieID}",
		h.removeMovie,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/collections/{collectionID}/shows/{showID}",
		h.addShow,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.DELETEI("/api/v1/collections/{collectionID}/shows/{showID}",
		h.removeShow,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *CollectionHandler) listCollections(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	collections, err := h.collectionService.ListByUserID(r.Context(), userID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read collections"})
		return
	}

	result := make([]CollectionInfo, 0, len(collections))
	for _, collection := range collections {
		result = append(result, toCollectionInfo(&collection))
	}

	respond.JSON(w, http.StatusOK, CollectionListResponse{Collections: result})
}

func (h *CollectionHandler) createCollection(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	requestBody, err := decodeCreateCollectionRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	collection, err := h.collectionService.Create(r.Context(), userID, service.CreateCollectionInput{Name: requestBody.Name})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create collection"})
		}
		return
	}

	respond.JSON(w, http.StatusCreated, toCollectionInfo(collection))
}

func (h *CollectionHandler) getCollection(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	collectionID := r.PathValue("collectionID")
	collection, err := h.collectionService.GetByID(r.Context(), userID, collectionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection id"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read collection"})
		}
		return
	}
	if collection == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		return
	}

	movies, err := h.collectionService.ListMovies(r.Context(), userID, collectionID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read collection movies"})
		return
	}

	shows, err := h.collectionService.ListShows(r.Context(), userID, collectionID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read collection shows"})
		return
	}

	movieInfos := make([]MovieInfo, 0, len(movies))
	for _, movie := range movies {
		movieInfos = append(movieInfos, toMovieInfo(&movie))
	}

	showInfos := make([]ShowInfo, 0, len(shows))
	for _, show := range shows {
		showInfos = append(showInfos, toShowInfo(&show))
	}

	respond.JSON(w, http.StatusOK, CollectionDetailResponse{
		Collection: toCollectionInfo(collection),
		Movies:     movieInfos,
		Shows:      showInfos,
	})
}

func (h *CollectionHandler) updateCollection(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	requestBody, err := decodeUpdateCollectionRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	collection, err := h.collectionService.UpdateName(r.Context(), userID, r.PathValue("collectionID"), service.UpdateCollectionInput{Name: requestBody.Name})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, toCollectionInfo(collection))
}

func (h *CollectionHandler) deleteCollection(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.collectionService.DeleteByID(r.Context(), userID, r.PathValue("collectionID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection id"})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CollectionHandler) addMovie(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.collectionService.AddMovie(r.Context(), userID, r.PathValue("collectionID"), r.PathValue("movieID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		case errors.Is(err, repository.ErrCollectionMovieAlreadyExists):
			respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to add movie to collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CollectionHandler) removeMovie(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.collectionService.RemoveMovie(r.Context(), userID, r.PathValue("collectionID"), r.PathValue("movieID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to remove movie from collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CollectionHandler) addShow(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.collectionService.AddShow(r.Context(), userID, r.PathValue("collectionID"), r.PathValue("showID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		case errors.Is(err, repository.ErrCollectionShowAlreadyExists):
			respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to add show to collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CollectionHandler) removeShow(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.collectionService.RemoveShow(r.Context(), userID, r.PathValue("collectionID"), r.PathValue("showID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCollectionInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid collection input"})
		case errors.Is(err, repository.ErrCollectionNotFound):
			respond.JSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to remove show from collection"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

func decodeCreateCollectionRequest(r *http.Request) (*CreateCollectionRequest, error) {
	defer r.Body.Close()

	var requestBody CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}

func decodeUpdateCollectionRequest(r *http.Request) (*UpdateCollectionRequest, error) {
	defer r.Body.Close()

	var requestBody UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}
