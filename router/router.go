package router

import (
	"go-api-portfolio/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/auth/login", middleware.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/welcome", middleware.Welcome).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/article/{id}", middleware.GetArticle).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/articles", middleware.GetAllArticles).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newarticle", middleware.CreateArticle).Methods("POST", "OPTIONS")
	//router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletearticle/{id}", middleware.DeleteArticle).Methods("DELETE", "OPTIONS")

	return router
}
