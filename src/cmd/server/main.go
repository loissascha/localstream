package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/database"
	"github.com/loissascha/localstream/internal/handler"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/provider/tmdb"
	"github.com/loissascha/localstream/internal/provider/tvmaze"
	repopostgres "github.com/loissascha/localstream/internal/repository/postgres"
	"github.com/loissascha/localstream/internal/service"
	backgroundservice "github.com/loissascha/localstream/internal/service/background"
)

func main() {
	godotenv.Load()

	logger.Config.ShowDebug(true)

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not defined. Make sure there is a .env file or a environment variable set!")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
		logger.Warning(nil, "No APP_ENV found. Setting 'development'")
	}

	listenAddr := fmt.Sprintf(":%v", port)
	if env == "development" {
		go func() {
			const pprofAddr = "127.0.0.1:6060"
			logger.Info(nil, "pprof listening on {addr}", pprofAddr)
			if err := http.ListenAndServe(pprofAddr, nil); err != nil {
				logger.Error(err, "pprof server stopped")
			}
		}()
	}

	s, err := server.NewServer(
		server.SetExportTypesLocation("frontend/src/lib/types/export_types.ts"),
		server.EnableExportTypes(env == "development"),
	)
	if err != nil {
		panic(err)
	}

	db, err := database.OpenPostgresFromEnv()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrationsDir := "./migrations"

	if err := database.RunMigrations(db.DB, migrationsDir); err != nil {
		panic(err)
	}

	// metadata providers
	tvMazeProvider := tvmaze.NewTVMazeProvider()
	tmdbProvider := tmdb.NewTMDBProvider()

	// repositories
	userRepo := repopostgres.NewUserRepository(db)
	libraryRepo := repopostgres.NewLibraryRepository(db)
	showRepo := repopostgres.NewShowRepository(db)
	showMetaRepo := repopostgres.NewShowMetadataRepository(db)
	seasonMetaRepo := repopostgres.NewSeasonMetadataRepository(db)
	episodeMetaRepo := repopostgres.NewEpisodeMetadataRepository(db)
	seasonRepo := repopostgres.NewSeasonRepository(db)
	episodeRepo := repopostgres.NewEpisodeRepository(db)
	userWatchstateRepo := repopostgres.NewUserWatchstateRepository(db)
	userMovieWatchstateRepo := repopostgres.NewUserMovieWatchstateRepository(db)
	movieRepo := repopostgres.NewMovieRepository(db)
	movieMetaRepo := repopostgres.NewMovieMetadataRepository(db)
	collectionRepo := repopostgres.NewCollectionRepository(db)

	// services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	libService := service.NewLibraryService(libraryRepo)
	showSerivce := service.NewShowService(showRepo)
	seasonService := service.NewSeasonService(seasonRepo)
	episodeService := service.NewEpisodeService(episodeRepo, seasonRepo)
	userWatchstateService := service.NewUserWatchstateService(userWatchstateRepo)
	userMovieWatchstateServiced := service.NewUserMovieWatchstateService(userMovieWatchstateRepo)
	movieService := service.NewMovieService(movieRepo)
	showMetaService := service.NewShowMetadataService(showMetaRepo, showRepo, tvMazeProvider)
	seasonMetaService := service.NewSeasonMetadataService(seasonMetaRepo)
	episodeMetaService := service.NewEpisodeMetadataService(episodeMetaRepo)
	movieMetaService := service.NewMovieMetadataService(movieService, movieMetaRepo, movieRepo, tmdbProvider)
	collectionService := service.NewCollectionService(collectionRepo)

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// handler
	authH := handler.NewAuthHandler(s, authService, authMiddleware)
	libH := handler.NewLibraryHandler(s, authMiddleware, libService)
	showH := handler.NewShowHandler(s, authMiddleware, showSerivce)
	seasonH := handler.NewSeasonHandler(s, authMiddleware, seasonService)
	episodeH := handler.NewEpisodeHandler(s, authMiddleware, episodeService, userWatchstateService)
	userWatchstateH := handler.NewUserWatchstateHandler(s, authMiddleware, userWatchstateService, showSerivce, seasonService, episodeService)
	userMovieWatchstateH := handler.NewUserMovieWatchstateHandler(s, authMiddleware, userMovieWatchstateServiced, movieService)
	movieH := handler.NewMovieHandler(s, authMiddleware, movieService)
	showMetaH := handler.NewShowMetadataHandler(s, authMiddleware, showMetaService)
	seasonMetaH := handler.NewSeasonMetadataHandler(s, authMiddleware, seasonMetaService)
	episodeMetaH := handler.NewEpisodeMetadataHandler(s, authMiddleware, episodeMetaService)
	movieMetaH := handler.NewMovieMetadataHandler(s, authMiddleware, movieMetaService)
	searchH := handler.NewSearchHandler(s, authMiddleware, showSerivce, movieService)
	collectionH := handler.NewCollectionHandler(s, authMiddleware, collectionService)

	// register routes
	authH.RegisterHandlers()
	libH.RegisterHandlers()
	showH.RegisterRoutes()
	seasonH.RegisterRoutes()
	episodeH.RegisterRoutes()
	userWatchstateH.RegisterRoutes()
	movieH.RegisterRoutes()
	userMovieWatchstateH.RegisterHandlers()
	showMetaH.RegisterRoutes()
	seasonMetaH.RegisterRoutes()
	episodeMetaH.RegisterRoutes()
	movieMetaH.RegisterRoutes()
	searchH.RegisterRoutes()
	collectionH.RegisterRoutes()

	frontendBuildDir := os.Getenv("FRONTEND_APP_DIR")
	frontendFileServer := http.FileServer(http.Dir(frontendBuildDir))
	s.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		requestPath := strings.TrimPrefix(r.URL.Path, "/")
		if requestPath != "" {
			filePath := filepath.Join(frontendBuildDir, filepath.Clean(requestPath))
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				frontendFileServer.ServeHTTP(w, r)
				return
			}

			if filepath.Ext(requestPath) != "" {
				http.NotFound(w, r)
				return
			}
		}

		http.ServeFile(w, r, filepath.Join(frontendBuildDir, "index.html"))
	})

	libraryCataloguer := backgroundservice.NewLibraryCataloguer(libService, movieMetaService, showRepo, seasonRepo, episodeRepo, movieRepo, tvMazeProvider, tmdbProvider, showMetaRepo, movieMetaRepo, seasonMetaRepo, episodeMetaRepo)
	libraryCataloguer.RunBackground()

	libraryUncataloguer := backgroundservice.NewLibraryUncataloguer(showRepo, seasonRepo, episodeRepo, movieRepo)
	libraryUncataloguer.RunBackground()

	logger.Info(nil, "Server starting at {addr}", listenAddr)
	err = s.Serve(listenAddr)
	if err != nil {
		logger.Error(err, "Server failed to start...")
	}
}
