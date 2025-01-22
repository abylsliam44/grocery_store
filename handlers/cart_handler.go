package handlers

import (
	"encoding/json"
	"final_project/database"
	"final_project/middlewares"
	"final_project/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	// Парсинг тела запроса
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.ProductID <= 0 || requestBody.Quantity <= 0 {
		log.Println("Invalid request body:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Получение user_id из сессии
	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Проверка наличия товара
	var product models.Product
	err = database.DB.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", requestBody.ProductID).
		Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		log.Println("Product not found:", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if product.Stock < requestBody.Quantity {
		log.Println("Not enough stock")
		http.Error(w, "Not enough stock available", http.StatusBadRequest)
		return
	}

	// Добавление в корзину
	_, err = database.DB.Exec(`
        INSERT INTO cart (user_id, product_id, quantity)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, product_id) DO UPDATE
        SET quantity = cart.quantity + $3
    `, userID, requestBody.ProductID, requestBody.Quantity)
	if err != nil {
		log.Println("Failed to add to cart:", err)
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	log.Printf("Product %d added to cart by user %d", requestBody.ProductID, userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product added to cart"))
}

func CartHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middlewares.Store.Get(r, "session-name")
	userID := session.Values["user_id"].(int)

	rows, err := database.DB.Query(`
        SELECT p.id, p.name, p.price, c.quantity, (p.price * c.quantity) AS total
        FROM cart c
        INNER JOIN products p ON c.product_id = p.id
        WHERE c.user_id = $1
    `, userID)
	if err != nil {
		log.Println("Error fetching cart:", err)
		http.Error(w, "Unable to fetch cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Price, &item.Quantity, &item.Total); err != nil {
			log.Println("Error scanning cart item:", err)
			http.Error(w, "Unable to fetch cart items", http.StatusInternalServerError)
			return
		}
		cartItems = append(cartItems, item)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/cart.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"CartItems": cartItems,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

func DeleteFromCartHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение product_id из формы
	r.ParseForm()
	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil || productID <= 0 {
		log.Println("Invalid Product ID:", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Получение user_id из сессии
	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Удаление товара из корзины
	_, err = database.DB.Exec("DELETE FROM cart WHERE user_id = $1 AND product_id = $2", userID, productID)
	if err != nil {
		log.Println("Failed to remove item from cart:", err)
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		return
	}

	log.Printf("Product %d removed from cart by user %d", productID, userID)
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
