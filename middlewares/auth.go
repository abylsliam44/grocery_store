package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("super-secret-key"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		userID, ok := session.Values["user_id"].(int)
		if !ok || userID <= 0 {
			
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request, userID int) {
	session, _ := Store.Get(r, "session-name")
	session.Values["user_id"] = userID
	session.Save(r, w)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	delete(session.Values, "user_id")
	session.Save(r, w)
}
