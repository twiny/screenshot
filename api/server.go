package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Server
type Server struct {
	srv    *http.Server
	router *chi.Mux
	app    *App
}

// NewHTTPServer
func NewHTTPServer(app *App) *Server {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)

		json.NewEncoder(w).Encode(struct {
			Status  int         `json:"Status"`
			Payload interface{} `json:"Payload"`
		}{
			Status:  404,
			Payload: "not found",
		})
	})

	// basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           30, // Maximum value not ignored by any of major browsers
	})

	// basic middleware
	r.Use(
		middleware.Compress(5, "gzip"),
		middleware.Timeout(15*time.Second),
		middleware.AllowContentType("application/json", "text/plain"),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.DefaultLogger,
		middleware.StripSlashes,
		cors.Handler,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.RequestID,
	)

	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		srv:    srv,
		router: r,
		app:    app,
	}
}

// writer
func (mux *Server) writer(w http.ResponseWriter, status int, payload interface{}) interface{} {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	//
	return json.NewEncoder(w).Encode(struct {
		Status  int         `json:"Status"`
		Payload interface{} `json:"Payload"`
	}{
		Status:  status,
		Payload: payload,
	})
}

// ListenAndServe
func (mux *Server) ListenAndServe(port string) error {
	// middleware
	mux.router.Use(mux.limiter)

	// routes
	mux.router.Post("/capture", mux.takeScreenShot)
	mux.router.Get("/download/{uuid}", mux.findScreenShot)
	mux.router.Get("/stats", mux.stats)

	mux.srv.Addr = port

	// log
	log.Println("http server running on " + mux.srv.Addr)

	// http
	return mux.srv.ListenAndServe()
}
