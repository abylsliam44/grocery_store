package handlers

import (
	"final_project/database"
	"final_project/middlewares"
	"fmt"
	"net/http"
)

func AddReviewHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok || userID <= 0 {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	productID := r.FormValue("product_id")
	rating := r.FormValue("rating")
	review := r.FormValue("review")

	_, err := database.DB.Exec("INSERT INTO reviews (product_id, user_id, rating, review) VALUES ($1, $2, $3, $4)", productID, userID, rating, review)
	if err != nil {
		http.Error(w, "Unable to save review", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/products/%s", productID), http.StatusSeeOther)
}
