package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cooperstandard/NetZero/internal/auth"
	"github.com/cooperstandard/NetZero/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db          *database.Queries
	tokenSecret string
	adminKey    string
	platform    string
}

const basePath = "/api/v1"

func main() {
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
		db:          dbQueries,
		platform:    platform,
		tokenSecret: os.Getenv("TOKEN_SECRET"),
		adminKey:    os.Getenv("ADMIN_KEY"),
	}

	apiMux := http.NewServeMux()

	// TODO: eventually, dynamically create a slice of routes on startup based on env.platform and then register with a for each
	// type route struct {
	// 	pattern string
	// 	handler http.HandlerFunc
	// }

	// routes
	register(apiMux, formPath("POST", "/reset"), apiCfg.adminAuth(apiCfg.handleReset))

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

func formPath(method, route string) string {
	return fmt.Sprintf("%s %s%s", method, basePath, route)
}

func (cfg apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Header.Get("Authorization"))
	w.WriteHeader(204)
}

func (cfg apiConfig) adminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil || token != cfg.adminKey {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (cfg apiConfig) userAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return next
}
