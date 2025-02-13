package handlers

import (
	"final_project/database"
	"final_project/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, content, created_at FROM blog_posts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Unable to fetch blog posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.BlogPost
	for rows.Next() {
		var post models.BlogPost
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Unable to parse blog data", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	fmt.Printf("Number of posts: %d\n", len(posts))
	tmpl, err := template.ParseFiles("templates/base.html", "templates/blog.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Posts": posts,
	})
}

// хэндлер для конкр поста

func BlogPostHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	postID := vars["id"]

	var post models.BlogPost
	err := database.DB.QueryRow("SELECT id, title, content, created_at FROM blog_posts WHERE id = $1", postID).
		Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/blog_post.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Post": post,
	})
}
