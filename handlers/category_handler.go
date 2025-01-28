package handlers

import (
	"final_project/database"
	"final_project/models"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Функция для отображения продуктов по категории
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	var categoryName string
	err := database.DB.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID).Scan(&categoryName)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(`
        SELECT id, name, price, stock, image_url
        FROM products
        WHERE category_id = $1
    `, categoryID)
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

	tmpl, err := template.ParseFiles("templates/base.html", "templates/category.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"CategoryName": categoryName,
		"Products":     products,
	}

	tmpl.Execute(w, data)
}

// Функция для отображения всех категорий
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Unable to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			http.Error(w, "Unable to parse category data", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/categories.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Categories": categories,
	}

	tmpl.Execute(w, data)
}
