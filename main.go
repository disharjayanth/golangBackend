package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/disharjayanth/golangBackend/data"
)

var temp *template.Template
var err error

func init() {
	temp, err = template.ParseGlob("template/*.html")
	if err != nil {
		log.Println("Failed to parse template files.", err)
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/signup", signUpPage)
	http.HandleFunc("/signin", signInPage)
	http.HandleFunc("/movie", moviePage)

	// handles css
	http.Handle("/stylesheet/", http.StripPrefix("/stylesheet", http.FileServer(http.Dir("template/stylesheet/"))))

	log.Println("Server serving @PORT: 3000")
	http.ListenAndServe(":3000", nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	temp.ExecuteTemplate(w, "mainPage.html", nil)
}

func signUpPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(200)
		temp.ExecuteTemplate(w, "signUp.html", nil)
	case "POST":
		name := r.FormValue("name")
		password := r.FormValue("password")
		fmt.Println("name password:", name, password)
		user := data.User{
			Name:     name,
			Password: password,
		}
		if user.Store() {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "signUp.html", "User already registered.")
		}
	}
}

func signInPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(200)
		temp.ExecuteTemplate(w, "signIn.html", nil)
	case "POST":
		name := r.FormValue("name")
		password := r.FormValue("password")
		fmt.Println("name password:", name, password)
		user := data.User{
			Name:     name,
			Password: password,
		}
		if user.Auth() {
			http.Redirect(w, r, "http://www.omdbapi.com/", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "signIn.html", "Username or password not correct!")
		}
	}
}

func moviePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Movie Page"))
}
