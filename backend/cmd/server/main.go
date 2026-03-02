package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/database"
	"github.com/loissascha/localstream/internal/handler"
	repopostgres "github.com/loissascha/localstream/internal/repository/postgres"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not defined. Make sure there is a .env file or a environment variable set!")
	}

	s, err := server.NewServer(
		server.SetExportTypesLocation("../export_types.ts"),
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

	_ = repopostgres.NewUserRepository(db)

	// handler
	videoH := handler.NewVideoHandler(s)

	// register routes
	videoH.RegisterHandlers()

	// fs := http.FileServer(http.Dir("./static"))
	// s.GetMux().Handle("/static/", http.StripPrefix("/static/", fs))

	logger.Info(nil, "Server starting at port: {port}", port)
	err = s.Serve(fmt.Sprintf(":%v", port))
	if err != nil {
		logger.Error(err, "Server failed to start...")
	}
}
