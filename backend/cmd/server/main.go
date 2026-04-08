package main

import (
	"fmt"
	"os"

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

	s, err := server.NewServer(
		server.SetExportTypesLocation("../frontend/src/lib/types/export_types.ts"),
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

	// repositories
	userRepo := repopostgres.NewUserRepository(db)
	libraryRepo := repopostgres.NewLibraryRepository(db)
	showRepo := repopostgres.NewShowRepository(db)
	showMetaRepo := repopostgres.NewShowMetadataRepository(db)
	seasonRepo := repopostgres.NewSeasonRepository(db)
	episodeRepo := repopostgres.NewEpisodeRepository(db)
	userWatchstateRepo := repopostgres.NewUserWatchstateRepository(db)
	userMovieWatchstateRepo := repopostgres.NewUserMovieWatchstateRepository(db)
	movieRepo := repopostgres.NewMovieRepository(db)
	movieMetaRepo := repopostgres.NewMovieMetadataRepository(db)

	// services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	libService := service.NewLibraryService(libraryRepo)
	showSerivce := service.NewShowService(showRepo)
	seasonService := service.NewSeasonService(seasonRepo)
	episodeService := service.NewEpisodeService(episodeRepo)
	userWatchstateService := service.NewUserWatchstateService(userWatchstateRepo)
	userMovieWatchstateServiced := service.NewUserMovieWatchstateService(userMovieWatchstateRepo)
	movieService := service.NewMovieService(movieRepo)
	showMetaService := service.NewShowMetadataService(showMetaRepo, showRepo)

	// middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// handler
	authH := handler.NewAuthHandler(s, authService, authMiddleware)
	videoH := handler.NewVideoHandler(s, authMiddleware)
	libH := handler.NewLibraryHandler(s, authMiddleware, libService)
	showH := handler.NewShowHandler(s, authMiddleware, showSerivce)
	seasonH := handler.NewSeasonHandler(s, authMiddleware, seasonService)
	episodeH := handler.NewEpisodeHandler(s, authMiddleware, episodeService, userWatchstateService)
	userWatchstateH := handler.NewUserWatchstateHandler(s, authMiddleware, userWatchstateService, showSerivce, seasonService, episodeService)
	userMovieWatchstateH := handler.NewUserMovieWatchstateHandler(s, authMiddleware, userMovieWatchstateServiced, movieService)
	movieH := handler.NewMovieHandler(s, authMiddleware, movieService)
	showMetaH := handler.NewShowMetadataHandler(s, authMiddleware, showMetaService)

	// register routes
	authH.RegisterHandlers()
	videoH.RegisterHandlers()
	libH.RegisterHandlers()
	showH.RegisterRoutes()
	seasonH.RegisterRoutes()
	episodeH.RegisterRoutes()
	userWatchstateH.RegisterRoutes()
	movieH.RegisterRoutes()
	userMovieWatchstateH.RegisterHandlers()
	showMetaH.RegisterRoutes()

	// metadata providers
	tvMazeProvider := tvmaze.NewTVMazeProvider()
	tmdbProvider := tmdb.NewTMDBProvider()

	libraryCataloguer := backgroundservice.NewLibraryCataloguer(libService, showRepo, seasonRepo, episodeRepo, movieRepo, tvMazeProvider, tmdbProvider, showMetaRepo, movieMetaRepo)
	libraryCataloguer.RunBackground()

	logger.Info(nil, "Server starting at port: {port}", port)
	err = s.Serve(fmt.Sprintf(":%v", port))
	if err != nil {
		logger.Error(err, "Server failed to start...")
	}
}
