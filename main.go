package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disharjayanth/golangBackend/data"
	"github.com/gorilla/sessions"
)

var temp *template.Template
var err error

var (
	// key must be 16, 24 or 32 bytes long
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type album struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
}

func init() {
	temp, err = template.ParseGlob("template/*.html")
	if err != nil {
		log.Println("Failed to parse template files.", err)
	}
}

func main() {
	// if go code crashes, it prints file name and also line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := http.Server{
		// for deployment add os.Getenv("PORT")
		Addr: ":" + os.Getenv("PORT"),
		// (for local development)
		// Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/signup", signUpPage)
	http.HandleFunc("/signin", signInPage)
	http.HandleFunc("/movie", moviePage)

	// handles css
	http.Handle("/stylesheet/", http.StripPrefix("/stylesheet", http.FileServer(http.Dir("template/stylesheet/"))))

	log.Println("Server serving @PORT: 3000")

	server.ListenAndServe()
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
			session, err := store.Get(r, "cookie-name")
			if err != nil {
				log.Println("Error creating session")
			}
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, "/movie", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "signIn.html", "Username or password incorrect!")
		}
	}
}

func moviePage(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		fmt.Println("Error accessing cookie info", err)
	}

	if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
		// temp.ExecuteTemplate(w, "signIn.html", "Please Sign In!")
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	} else {
		url := "https://jsonplaceholder.typicode.com/albums"

		req, _ := http.NewRequest("GET", url, nil)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making request to client:", err)
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var albums []album
		err = json.Unmarshal(body, &albums)
		if err != nil {
			log.Println("Error:", err)
		}

		temp.ExecuteTemplate(w, "movie.html", albums)
	}
}
