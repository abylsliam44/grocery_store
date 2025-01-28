package handlers

import (
	"final_project/database"
	"final_project/middlewares"
	"final_project/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры сортировки и поиска
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	searchQuery := r.URL.Query().Get("search_query")

	if sortOrder == "" {
		sortOrder = "asc" // по умолчанию сортируем по возрастанию
	}
	if sortBy == "" {
		sortBy = "price" // по умолчанию сортируем по цене
	}

	// Строим SQL запрос с учетом сортировки и поиска
	query := fmt.Sprintf(`
        SELECT id, name, price, stock, image_url
        FROM products
        WHERE name ILIKE $1 OR category_id::text ILIKE $2
        ORDER BY %s %s
    `, sortBy, sortOrder)

	rows, err := database.DB.Query(query, "%"+searchQuery+"%", "%"+searchQuery+"%")
	if err != nil {
		http.Error(w, "Unable to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.ImageURL); err != nil {
			http.Error(w, "Unable to parse product data", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	// Получаем информацию о том, авторизован ли пользователь
	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	isAuthenticated := false
	if ok && userID > 0 {
		isAuthenticated = true
	}

	// Рендерим шаблон
	tmpl, err := template.ParseFiles("templates/base.html", "templates/products.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Products":        products,
		"SortBy":          sortBy,
		"SortOrder":       sortOrder,
		"SearchQuery":     searchQuery,
		"IsAuthenticated": isAuthenticated,
	}

	tmpl.Execute(w, data)
}

func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр id из URL
	vars := mux.Vars(r)
	productID := vars["id"]

	// Получаем данные о товаре
	var product models.Product
	err := database.DB.QueryRow("SELECT id, name, price, image_url FROM products WHERE id = $1", productID).
		Scan(&product.ID, &product.Name, &product.Price, &product.ImageURL)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Получаем отзывы для товара
	rows, err := database.DB.Query("SELECT r.rating, r.review, u.name FROM reviews r JOIN users u ON r.user_id = u.id WHERE r.product_id = $1 ORDER BY r.created_at DESC", productID)
	if err != nil {
		http.Error(w, "Unable to fetch reviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.Rating, &review.Review, &review.User.Name); err != nil {
			http.Error(w, "Unable to parse review data", http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, review)
	}

	// Проверка, авторизован ли пользователь
	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	isAuthenticated := false
	if ok && userID > 0 {
		isAuthenticated = true
	}

	// Рендерим шаблон с товаром и отзывами
	tmpl, err := template.ParseFiles("templates/base.html", "templates/product_detail.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон
	data := map[string]interface{}{
		"Product":         product,
		"Reviews":         reviews,
		"IsAuthenticated": isAuthenticated, // Добавляем флаг для отображения формы отзыва
	}

	tmpl.Execute(w, data)
}
