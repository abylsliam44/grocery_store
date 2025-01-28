package handlers

import (
	"final_project/database"
	"final_project/middlewares"
	"final_project/models"
	"html/template"
	"log"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из сессии
	session, _ := middlewares.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	isAuthenticated := false

	if ok && userID > 0 {
		isAuthenticated = true
	} else {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Если метод GET, то просто выводим профиль
	if r.Method == http.MethodGet {
		var user models.User
		err := database.DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Println("Error fetching user data:", err)
			http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/base.html", "templates/profile.html")
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"User":            user,
			"IsAuthenticated": isAuthenticated,
		}

		tmpl.Execute(w, data)
		return
	}

	// Если метод POST, то обновляем информацию
	if r.Method == http.MethodPost {
		// Получаем новые данные из формы
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Обновляем данные пользователя в базе
		_, err := database.DB.Exec("UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4", name, email, password, userID)
		if err != nil {
			log.Println("Error updating user data:", err)
			http.Error(w, "Unable to update user data", http.StatusInternalServerError)
			return
		}

		// Перенаправляем на страницу профиля после обновления
		http.Redirect(w, r, "/users/profile", http.StatusSeeOther)
		return
	}
}
