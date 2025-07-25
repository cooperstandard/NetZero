package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/cooperstandard/NetZero/internal/routes"
	"github.com/cooperstandard/NetZero/internal/util"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func main() {
	const port = "8080"
	const basePath = "/api/v1"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("platform must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := routes.ApiConfig{
		DB:          dbQueries,
		Platform:    platform,
		TokenSecret: os.Getenv("TOKEN_SECRET"),
		AdminKey:    os.Getenv("ADMIN_KEY"),
	}

	apiMux := http.NewServeMux()

	// TODO: eventually, dynamically create a slice of routes on startup based on env.platform and then register with a for each
	// type route struct {
	// 	pattern string
	// 	handler http.HandlerFunc
	// }

	// routes
	register(apiMux, util.FormPath("POST", "/reset", basePath), apiCfg.AdminAuth(apiCfg.HandleReset))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: apiMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func register(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, handler)
}
