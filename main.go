package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Actor struct {
	ID          int
	Name        string
	Birthdate   string
	Nationality string
}

type Movie struct {
	ID           int
	Title        string
	Release_date string
	Director     string
}

type MovieActor struct {
	Movie_id int
	Actor_id int
}

func main() {
	db, err := sql.Open("sqlite3", "actors.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		birthdate := r.FormValue("birthdate")
		nationality := r.FormValue("nationality")

		stmt, err := db.Prepare("INSERT INTO actors (name, birthdate, nationality) VALUES (?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, birthdate, nationality)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("static/submit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, map[string]interface{}{"name": name})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("static/form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
