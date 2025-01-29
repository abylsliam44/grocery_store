package main

import (
	"log"
	"net/http"

	"final_project/database"
	"final_project/handlers"
	"final_project/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к базе данных
	database.Connect()

	// Инициализация маршрутов
	r := mux.NewRouter()

	// Открытые маршруты
	r.HandleFunc("/users/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/users/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/test", handlers.TestHandler)
	r.HandleFunc("/categories", handlers.CategoriesHandler).Methods("GET")
	r.HandleFunc("/products/category/{id}", handlers.CategoryHandler).Methods("GET")
	r.HandleFunc("/add-review", handlers.AddReviewHandler).Methods("GET", "POST")
	r.HandleFunc("/products/{id}", handlers.ProductDetailHandler).Methods("GET")
	r.HandleFunc("/blog", handlers.BlogHandler).Methods("GET")
	r.HandleFunc("/blog/{id:[0-9]+}", handlers.BlogPostHandler).Methods("GET")

	// Защищённые маршруты
	r.Handle("/products", middlewares.AuthMiddleware(http.HandlerFunc(handlers.ProductsHandler))).Methods("GET")
	r.Handle("/cart", middlewares.AuthMiddleware(http.HandlerFunc(handlers.CartHandler))).Methods("GET")
	r.Handle("/cart", middlewares.AuthMiddleware(http.HandlerFunc(handlers.AddToCartHandler))).Methods("POST")
	r.Handle("/cart/delete", middlewares.AuthMiddleware(http.HandlerFunc(handlers.DeleteFromCartHandler))).Methods("POST")

	r.Handle("/users/logout", middlewares.AuthMiddleware(http.HandlerFunc(handlers.LogoutHandler))).Methods("GET")

	r.Handle("/users/profile", middlewares.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler))).Methods("GET", "POST")

	// Статические файлы
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
