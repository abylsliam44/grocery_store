package handlers

import (
	"html/template"
	"net/http"

	"final_project/database"
	"final_project/middlewares"
	"final_project/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, password)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user models.User
	err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	session, _ := middlewares.Store.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	middlewares.LogoutUser(w, r)
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
