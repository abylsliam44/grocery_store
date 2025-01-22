package handlers

import (
	"fmt"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/test.html")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		productID := r.FormValue("product_id")
		quantity := r.FormValue("quantity")

		fmt.Fprintf(w, "Received ProductID: %s\n", productID)
		fmt.Fprintf(w, "Received Quantity: %s\n", quantity)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
