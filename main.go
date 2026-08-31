package main

import (
	"blogging-api/authenticator"
	"blogging-api/blogs"
	"blogging-api/db"
	"blogging-api/token"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error to load `.env` file")
	}
	db_config := db.SetConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		db_config.Host,
		db_config.Port,
		db_config.User,
		db_config.Password,
		db_config.Name,
	))
	//========TABLE INITIALIZATION========
	_, err_table_bloggers := db.Exec(`CREATE TABLE IF NOT EXISTS bloggers( 
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		nickname VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	)`)
	if err_table_bloggers != nil {
		panic("Error to create table `bloggers`: " + err_table_bloggers.Error())
	}
	_, err_table_blogs := db.Exec(`CREATE TABLE IF NOT EXISTS blogs(
		id SERIAL PRIMARY KEY,
		blogger_id INT,
		title VARCHAR(255) NOT NULL,
		content VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL,
		tags TEXT[],

		CONSTRAINT fk_blogs_bloggers
			FOREIGN KEY (blogger_id)
			REFERENCES bloggers(id)
	)`)
	if err_table_blogs != nil {
		panic("Error to create table `blogs`: " + err_table_blogs.Error())
	}
	//====================================
	if err != nil {
		panic(fmt.Sprintf("Error to open DB: %s", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Error to ping DB: %s", err))
	}
	token := token.Token{
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	}
	authenticator := authenticator.Authenticator{
		DB:    db,
		Token: &token,
	}

	blogs := blogs.Blogs{
		DB:    db,
		Token: &token,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authenticator.Registration)
	mux.HandleFunc("POST /login", authenticator.Login)

	mux.HandleFunc("POST /posts", blogs.Create)
	mux.HandleFunc("PUT /posts/{id}", blogs.Update)
	mux.HandleFunc("DELETE /posts/{id}", blogs.Delete)

	mux.HandleFunc("GET /posts", blogs.Get)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic("Error to set connection: " + err.Error())
	}
}
