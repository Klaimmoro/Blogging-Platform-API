package authenticator

import (
	dberrors "blogging-api/db/db_errors"
	"blogging-api/token"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Authenticator struct {
	DB    *sql.DB
	Token *token.Token
}

func (a *Authenticator) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var reg_req RegistrationRequest
	err_parse := json.NewDecoder(r.Body).Decode(&reg_req)
	if err_parse != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	hashed_password, err_hashing := bcrypt.GenerateFromPassword([]byte(reg_req.Password), bcrypt.DefaultCost)
	if err_hashing != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var blogger_id int
	insert_row := a.DB.QueryRow("INSERT INTO bloggers (email, nickname, password) VALUES ($1, $2, $3) RETURNING id",
		reg_req.Email,
		reg_req.Nickname,
		hashed_password,
	).Scan(&blogger_id)
	if insert_row != nil {
		if dberrors.IsUnique(insert_row) {
			http.Error(w, "Email already registred", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	auth_token, err_token := a.Token.GenerateToken(blogger_id)
	if err_token != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": auth_token})
}

func (a *Authenticator) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var log_req LoginRequest
	err_parse := json.NewDecoder(r.Body).Decode(&log_req)
	if err_parse != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var blogger_id int
	var hashed_password string
	select_err := a.DB.QueryRow(`SELECT id, password FROM bloggers WHERE email=$1`, log_req.Email).Scan(&blogger_id, &hashed_password)
	if select_err != nil {
		if errors.Is(select_err, sql.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		log.Println("Error to generate token: " + select_err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	check_pass := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(log_req.Password))
	if check_pass != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	auth_token, err_token := a.Token.GenerateToken(blogger_id)
	if err_token != nil {
		log.Println("Error to generate token: " + err_token.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"token": auth_token})
}
