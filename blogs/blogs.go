package blogs

import (
	"blogging-api/token"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type BlogRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type BlogResponse struct {
	ID        int      `json:"id"`
	BloggerID int      `json:"blogger_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
}

type Blogs struct {
	DB    *sql.DB
	Token *token.Token
}

func (b *Blogs) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	auth_data := r.Header.Get("Authorization")
	auth_token := strings.TrimPrefix(auth_data, "Bearer ")
	user_id, err_parse_token := b.Token.ParseToken(auth_token)
	if err_parse_token != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	defer r.Body.Close()
	var req_blog BlogRequest
	err_parse_req := json.NewDecoder(r.Body).Decode(&req_blog)
	if err_parse_req != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var blog_id int
	db_insert := b.DB.QueryRow(`INSERT INTO blogs (blogger_id, title,content,category,tags) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		user_id,
		req_blog.Title,
		req_blog.Content,
		req_blog.Category,
		req_blog.Tags,
	).Scan(&blog_id)
	if db_insert != nil {
		if errors.Is(db_insert, sql.ErrNoRows) {
			http.Error(w, "Unknow blogger", http.StatusConflict)
			return
		}
		log.Println("Error to inser blog: " + db_insert.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(BlogResponse{
		ID:        blog_id,
		BloggerID: user_id,
		Title:     req_blog.Title,
		Content:   req_blog.Content,
		Category:  req_blog.Category,
		Tags:      req_blog.Tags,
	})
}

func (b *Blogs) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	auth_data := r.Header.Get("Authorization")
	auth_token := strings.TrimPrefix(auth_data, "Bearer ")
	user_id, err_parse_token := b.Token.ParseToken(auth_token)
	if err_parse_token != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	defer r.Body.Close()
	var req_blog BlogRequest
	err_parse_req := json.NewDecoder(r.Body).Decode(&req_blog)
	if err_parse_req != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	blog_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}
	var owner_id int
	db_select := b.DB.QueryRow(`SELECT blogger_id FROM blogs WHERE id=$1`, blog_id).Scan(&owner_id)
	if db_select != nil {
		if errors.Is(db_select, sql.ErrNoRows) {
			http.Error(w, "Unknow task", http.StatusConflict)
			return
		}
		log.Println("Error to inser blog: " + db_select.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user_id != owner_id {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden"})
		return
	}
	_, err_update := b.DB.Exec(`UPDATE blogs SET title=$1, content=$2, category=$3, tags=$4 WHERE id=$5`,
		req_blog.Title,
		req_blog.Content,
		req_blog.Category,
		req_blog.Tags,
		user_id,
	)
	if err_update != nil {
		if errors.Is(err_update, sql.ErrNoRows) {
			http.Error(w, "Unknow blogger", http.StatusConflict)
			return
		}
		log.Println("Error to inser blog: " + err_update.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(BlogResponse{
		ID:       blog_id,
		Title:    req_blog.Title,
		Content:  req_blog.Content,
		Category: req_blog.Category,
		Tags:     req_blog.Tags,
	})
}

func (b *Blogs) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	auth_data := r.Header.Get("Authorization")
	auth_token := strings.TrimPrefix(auth_data, "Bearer ")
	user_id, err_parse_token := b.Token.ParseToken(auth_token)
	if err_parse_token != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	blog_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}
	var owner_id int
	db_select := b.DB.QueryRow(`SELECT blogger_id FROM blogs WHERE id=$1`, blog_id).Scan(&owner_id)
	if db_select != nil {
		if errors.Is(db_select, sql.ErrNoRows) {
			http.Error(w, "Unknow task", http.StatusConflict)
			return
		}
		log.Println("Error to inser blog: " + db_select.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user_id != owner_id {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden"})
		return
	}
	_, err_update := b.DB.Exec(`DELETE FROM blogs WHERE id=$1`,
		blog_id,
	)
	if err_update != nil {
		if errors.Is(err_update, sql.ErrNoRows) {
			http.Error(w, "Unknow blogger", http.StatusConflict)
			return
		}
		log.Println("Error to insert blog: " + err_update.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (b *Blogs) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	auth_data := r.Header.Get("Authorization")
	auth_token := strings.TrimPrefix(auth_data, "Bearer ")
	user_id, err_parse_token := b.Token.ParseToken(auth_token)
	if err_parse_token != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	rows, err := b.DB.Query(`SELECT * FROM blogs WHERE blogger_id=$1`, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Blogs are not found", http.StatusNotFound)
			return
		}
		log.Println("Error to find blogs: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var blogs []BlogResponse
	defer rows.Close()
	for rows.Next() {
		var blog BlogResponse
		if err := rows.Scan(&blog.ID, &user_id, &blog.Title, &blog.Content, &blog.Category, pq.Array(&blog.Tags)); err != nil {
			log.Println("Error to scan row: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error to scan row: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}
