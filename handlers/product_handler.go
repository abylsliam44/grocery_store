package handlers

import (
	"final_project/database"
	"final_project/models"
	"html/template"
	"net/http"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		http.Error(w, "Unable to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			http.Error(w, "Unable to parse product data", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/products.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Products": products,
	}

	tmpl.Execute(w, data)
}
