package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	_ "github.com/Brightwater/recipe-backend-go/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	cors := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
		},
	)
	r.Use(cors.Handler)

	r.Use(middleware.Logger) // add log middleware

	// Serve the index.html file as the root path
	// r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, path.Join("static", "index.html"))
	// })

	// Serve static files from the "public" directory
	// r.Handle("/_app/*", http.FileServer(http.Dir("static")))
	// r.Handle("/ui/*", http.FileServer(http.Dir("static")))
	// r.Handle("/fonts/*", http.FileServer(http.Dir("static")))

	r.Get("/", http.RedirectHandler("/docs/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/docs", http.RedirectHandler("/docs/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/docs/*", httpSwagger.Handler())

	r.Get("/helloworld", hello)
	r.Get("/testAuth", testAuth)

	return r
}

// @Summary Hello World
// @Produce json
// @Success 200 {string} string "ok"
// @Router /helloworld [get]
func hello(w http.ResponseWriter, r *http.Request) {

	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		// Handle the error here, you can log it or return an error response
		fmt.Println("Error writing response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// @Summary Test Authentication
// @Success 200 {string} string "ok"
// @Success 401
// @Param jwt query string true "JWT"
// @Router /testAuth [get]
func testAuth(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")

	err := VerifyTokenAndScope(jwt)
	if err != nil {
		// return auth error to client
		w.WriteHeader(http.StatusUnauthorized)
	}

	// return ok to client
	w.WriteHeader(http.StatusOK)
}
