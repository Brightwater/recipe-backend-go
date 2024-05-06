package main

import (
	"encoding/json"
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
	r.Get("/recipe/getAllRecipes", getAllRecipes)
	r.Post("/recipe/addRecipe/", addRecipe)

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
// @Failure 401
// @Param jwt query string true "JWT"
// @Router /testAuth [get]
func testAuth(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")

	_, err := VerifyTokenAndScope(jwt)
	if err != nil {
		// return auth error to client
		w.WriteHeader(http.StatusUnauthorized)
	}

	// return ok to client
	w.WriteHeader(http.StatusOK)
}

// @Summary Gets all recipes from the database
// @Success 200 {string} string "ok"
// @Failure 500 {string} server error
// @Router /recipe/getAllRecipes [get]
func getAllRecipes(w http.ResponseWriter, r *http.Request) {

	recipes, err := GetAllRecipes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// @Summary Adds a new recipe into the database
// @Accept json
// @Param jwt query string true "JWT"
// @Param recipe body Recipe true "Recipe data"
// @Success 200 {string} string "ok"
// @Failure 400 {string} bad request
// @Failure 401
// @Failure 500 {string} server error
// @Router /recipe/addRecipe [post]
func addRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("adding recipe")
	// param is a jwt token
	jwt := r.URL.Query().Get("jwt")

	username, err := VerifyTokenAndScope(jwt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// body is a recipe
	var recipe Recipe = Recipe{}

	err = json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Error: " + err.Error()))
		if err != nil {
			fmt.Println("Error writing response:", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
		return
	}

	recipe.Author = username

	err = AddRecipe(&recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error: " + err.Error()))
		if err != nil {
			fmt.Println("Error writing response:", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
