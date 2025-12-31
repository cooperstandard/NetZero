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
	"github.com/pressly/goose/v3"
)

const migrationsDir = "./sql/migrations"

func main() {
	const port = "8080"
	const basePath = "/api/v1"

	// TODO: update example env
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
	defer dbConn.Close()
	dbQueries := database.New(dbConn)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose: failed to set dialect: %v\n", err)
	}

	if os.Getenv("RELOAD_MIGRATIONS") == "true" && platform != "prod" {
		goose.DownTo(dbConn, migrationsDir, 0)
	} 

	if err := goose.Up(dbConn, migrationsDir); err != nil {
		log.Fatalf("goose: failed to apply migrations: %v\n", err)
	}
	if err := goose.Version(dbConn, migrationsDir); err != nil {
		log.Fatalf("unable to get version")
	}

	log.Println("connected to DB starting server")

	apiCfg := routes.APIConfig{
		DBConn:      dbConn,
		DB:          dbQueries,
		Platform:    platform,
		TokenSecret: os.Getenv("TOKEN_SECRET"),
		AdminKey:    os.Getenv("ADMIN_KEY"),
	}

	apiMux := http.NewServeMux()

	paths := make(map[string]http.HandlerFunc)

	// add routes
	paths[util.FormPath("GET", "/admin/users", basePath)] = apiCfg.AdminAuthMiddleware(apiCfg.HandleGetUsers)
	paths[util.FormPath("GET", "/health", basePath)] = routes.HandleHealth
	paths[util.FormPath("POST", "/login", basePath)] = apiCfg.HandleLogin
	paths[util.FormPath("POST", "/register", basePath)] = apiCfg.HandleRegister
	paths[util.FormPath("POST", "/token/refresh", basePath)] = apiCfg.HandleRefreshToken
	paths[util.FormPath("POST", "/groups", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleCreateGroup)
	paths[util.FormPath("GET", "/groups", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleGetGroups)
	paths[util.FormPath("GET", "/groups/all", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleGetAllGroups)
	paths[util.FormPath("POST", "/groups/join", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleJoinGroup)
	paths[util.FormPath("GET", "/groups/members/{groupID}", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleGetMembers)
	paths[util.FormPath("POST", "/transaction", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleCreateTransactions)
	paths[util.FormPath("GET", "/transactions", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleGetTransactions)
	paths[util.FormPath("GET", "/transactions/details", basePath)] = apiCfg.UserAuthMiddleware(apiCfg.HandleGetTransactionDetails)

	if apiCfg.Platform != "prod" {
		paths[util.FormPath("POST", "/admin/reset", basePath)] = apiCfg.AdminAuthMiddleware(apiCfg.HandleReset)
		paths[util.FormPath("POST", "/admin/migrate", basePath)] = apiCfg.AdminAuthMiddleware(apiCfg.HandleMigration)
	}

	// register routes
	for k, v := range paths {
		register(apiMux, k, v)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: apiMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func register(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, routes.LogMiddleware(handler))
}

/* TODO: send refresh token like this, not in json response. Update the refresh endpoint to pull cookie from that cookie

  // Create a new cookie
		cookie := &http.Cookie{
			Name:     "my_session_cookie",
			Value:    "some_secret_session_id",
			Expires:  time.Now().Add(24 * time.Hour), // Set cookie to expire in 24 hours
			HttpOnly: true,                           // This is the key for HTTP-only
			Secure:   true,                           // Recommended for production (only send over HTTPS)
			SameSite: http.SameSiteLax,               // Recommended for production
			Path:     "/",                            // Available across the entire domain
		}

		// Add the cookie to the response
		http.SetCookie(w, cookie)

		fmt.Fprintf(w, "HTTP-only cookie set!")
*/
