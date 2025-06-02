package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "./static/"
	const port = "8080"

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

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()
	apiMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	apiMux.HandleFunc("GET /healthz", handlerReadiness)
	apiMux.HandleFunc("POST /validate_chirp", handlerChirpsValidate)
	apiMux.HandleFunc("POST /users", apiCfg.handleUsers)
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	adminMux.HandleFunc("POST /reset", apiCfg.handlerReset)
	adminMux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
