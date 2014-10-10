package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	//"time"
	"os"
	"text/template"
)
import "encoding/json"

const sessionName string = "user-session"


var store = sessions.NewCookieStore([]byte("something-very-secret"))

var router = mux.NewRouter()

func init() {

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 1, // 1 hour
		HttpOnly: true,
	}

}

func main() {

	port := ":8080"

	fmt.Println("Start Projectile Backend")

	// Rest API
	router.HandleFunc("/api/users", RestUsersHandler)
	router.HandleFunc("/api/user/{user}", RestUserHandler)

	router.HandleFunc("/api/items", RestItemsHandler)

	router.HandleFunc("/welcome", LandingPageHandler)
	router.HandleFunc("/signIn", SignInHandler)

	router.HandleFunc("/dashboard", DashboardHandler)
	router.HandleFunc("/items", ItemsHandler)

	router.HandleFunc("/show", ShowHandler)
	router.HandleFunc("/signUp", SignUpHandler)
	router.HandleFunc("/logOut", LogOutHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.HandleFunc("/", SessionHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	http.Handle("/", router)

	fmt.Println("Listen on", port)
	http.ListenAndServe(port, nil)
}

func NotFoundHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	fmt.Println("not found: ", path)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func ShowHandler(res http.ResponseWriter, req *http.Request) {

	//fmt.Println("Fuuuuu")

}

func CheckIfSessionIsValid(res http.ResponseWriter, req *http.Request) bool {

	session, _ := store.Get(req, sessionName)
	if val, ok := session.Values["username"].(string); ok {
		// if val is a string
		switch val {
		case "":
			return false
		default:
			return true
		}
	} else {
		return false
	}
}

var dbPasswdJson []struct {
	Username  string `json:"Username"`
	Password  string `json:"Password"`
}

func userValid(username, password string) (bool) {
	configFile, err := os.Open("db/passwd.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
		return false
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&dbPasswdJson); err != nil {
		fmt.Println("parsing config file", err.Error())
		return false
	}

	// search hashed password
	hashedPassword := ""
	for _, entry := range dbPasswdJson {
		if entry.Username == username {
			hashedPassword = entry.Password
		}
	}
	if hashedPassword == "" {
		fmt.Println("User not in database or password not available")
		return false
	}

	return CheckPassword([]byte(password), []byte(hashedPassword))
}

//handler for signIn
func SignInHandler(res http.ResponseWriter, req *http.Request) {
	username, password := req.FormValue("username"), req.FormValue("password")

	if !userValid(username, password) {
		fmt.Printf("Password for user %s invalid\n", username)
		http.Redirect(res, req, "/welcome", http.StatusFound)
	}

	sessionNew, _ := store.Get(req, sessionName)

	sessionNew.Values["username"] = username
	err := sessionNew.Save(req, res)
	if err != nil {
		panic("session save error")
	}
	store.Save(req, res, sessionNew)

	http.Redirect(res, req, "/dashboard", http.StatusFound)
}

//handler for signUp
func SignUpHandler(res http.ResponseWriter, req *http.Request) {

}

//handler for logOut
func LogOutHandler(res http.ResponseWriter, req *http.Request) {
	sessionOld, err := store.Get(req, sessionName)

	fmt.Println("Session in logout")
	fmt.Println(sessionOld)
	if err = sessionOld.Save(req, res); err != nil {
		fmt.Println("Error saving session: %v", err)
	}

	sessionOld.Options.MaxAge = -1
	_ = sessionOld.Save(req, res)

	http.Redirect(res, req, "/welcome", http.StatusFound)
}

func loadPageTemplate(title string) (*template.Template, error) {
	filename := "page-templates/" + title + ".html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type Person struct {
	Name string
}


func LandingPageHandler(res http.ResponseWriter, req *http.Request) {
	p, err := loadPageTemplate("welcome")
	if err != nil {
		http.Error(res, "Failed to load welcome page", http.StatusInternalServerError)
	}
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}

func SessionHandler(res http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, sessionName)

	if val, ok := session.Values["username"].(string); ok {
		// if val is a string
		switch val {
		case "":
			http.Redirect(res, req, "/welcome", http.StatusFound)
		default:
			http.Redirect(res, req, "/dashboard", http.StatusFound)
		}
	} else {
		// if val is not a string type
		http.Redirect(res, req, "/welcome", http.StatusFound)
	}
}
