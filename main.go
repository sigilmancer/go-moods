package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
    "path/filepath"

    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
	"go-mood-tracker/handlers"
	"go-mood-tracker/storage"
)

// Define a clear response shape
type CallingResponse struct {
    Message   string `json:"message"` // Use 'message' instead of 'text' for clarity
    Timestamp int64  `json:"timestamp"`
}

func main() {
    // 1. Setup Dependencies
	store := storage.NewStore()
	h := &handlers.Handler{Store: store}

	// 2. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	// 3. Define API Routes (Grouped)
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", h.Hello)
		
		r.Get("/pulse", h.GetPulse)
		r.Post("/submit", h.SubmitMood)
		 r.Delete("/delete", h.DeleteMood) 
	})
	r.Handle("/metrics", promhttp.Handler())
	// 4. SPA/Static File Serving
	setupFileServer(r, "./dist")

	// 5. Configure Production-Ready Server
	port := getEnv("PORT", "5000")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,  // Protect against slowloris
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
    log.Fatal(srv.ListenAndServe())
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setupFileServer(r chi.Router, publicPath string) {
	// 1. Create a standard file server
    fs := http.FileServer(http.Dir(publicPath))

    // 2. Catch-all route for static assets AND SPA routing
    r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
        // Build the full path to the requested file
        fullPath := filepath.Join(publicPath, r.URL.Path)

        // Check if the file exists on disk
        _, err := os.Stat(fullPath)
        
        // If it doesn't exist (or is a directory), serve index.html 
        // This lets Svelte handle the routing (SPA mode)
        if os.IsNotExist(err) {
            http.ServeFile(w, r, filepath.Join(publicPath, "index.html"))
            return
        }

        // Otherwise, serve the actual file (JS, CSS, Images)
        fs.ServeHTTP(w, r)
    })
}