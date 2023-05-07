package main

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var client *redis.Client
var templates *template.Template
var store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGethandler).Methods("GET")
	r.HandleFunc("/", indexPosthandler).Methods("POST")

	r.HandleFunc("/login", loginGethandler).Methods("GET")
	r.HandleFunc("/login", loginPosthandler).Methods("POST")
	r.HandleFunc("/register", registerGethandler).Methods("GET")
	r.HandleFunc("/register", registerPosthandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
func indexGethandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)

}
func indexPosthandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/", 302)
}

func loginGethandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPosthandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Save(r, w)
}
func registerGethandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}
func registerPosthandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return
	}
	client.Set("user:"+username, hashedPassword, 0)
	http.Redirect(w, r, "/login", 302)
}
