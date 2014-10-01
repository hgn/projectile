package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	//"time"
	//"os"
	"text/template"
)

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

	// Rest API
	router.HandleFunc("/api/users", RestUsersHandler)
	router.HandleFunc("/api/user/{user}", RestUserHandler)

	router.HandleFunc("/welcome", WelcomeHandler)
	router.HandleFunc("/signIn", SignInHandler)
	router.HandleFunc("/signInP", SignInPHandler)

	router.HandleFunc("/dashboard", DashboardHandler)
	router.HandleFunc("/items", ItemsHandler)

	router.HandleFunc("/show", ShowHandler)
	router.HandleFunc("/signUp", SignUpHandler)
	router.HandleFunc("/logOut", LogOutHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.HandleFunc("/", SessionHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("not found\n")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ShowHandler(res http.ResponseWriter, req *http.Request) {

	//fmt.Println("Fuuuuu")

}

func CheckIfSessionIsValid(res http.ResponseWriter, req *http.Request) bool {

	session, _ := store.Get(req, "loginSession")
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

//handler for signIn
func SignInHandler(res http.ResponseWriter, req *http.Request) {

	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Printf("username: %s\n", username)
	fmt.Printf("password: %s\n", password)

	//store in session variable
	sessionNew, _ := store.Get(req, "loginSession")

	// Set some session values.
	sessionNew.Values["username"] = username
	//sessionNew.Options.Secure = true
	// 10 minutes
	//sessionNew.Options.MaxAge = 60 * 10

	// Save it.
	err := sessionNew.Save(req, res)
	if err != nil {
		// handle the error case
	}
	store.Save(req, res, sessionNew)

	fmt.Println("Session after logging:")
	fmt.Println(sessionNew)

	http.Redirect(res, req, "/dashboard", http.StatusFound)
}

//handler for signUp
func SignUpHandler(res http.ResponseWriter, req *http.Request) {

}

//handler for logOut
func LogOutHandler(res http.ResponseWriter, req *http.Request) {
	sessionOld, err := store.Get(req, "loginSession")

	fmt.Println("Session in logout")
	fmt.Println(sessionOld)
	if err = sessionOld.Save(req, res); err != nil {
		fmt.Println("Error saving session: %v", err)
	}

	sessionOld.Options.MaxAge = -1
	_ = sessionOld.Save(req, res)

	http.Redirect(res, req, "/welcome", http.StatusFound)
}

func loadPage(title string) (*template.Template, error) {
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


func SignInPHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter SignInPHandler")

	var ret = CheckIfSessionIsValid(res, req)
	if ret == false && req.URL.Path != "/signInP" {
		http.Redirect(res, req, "/signInP", http.StatusFound)
		return
	}

	p, err := loadPage("sigin")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}
	fmt.Println(p)
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}

func WelcomeHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter WelcomeHandler")

	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/

	p, err := loadPage("welcome")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}
	fmt.Println(p)
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}

func SessionHandler(res http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "loginSession")

	if val, ok := session.Values["username"].(string); ok {
		// if val is a string
		switch val {
		case "":
			http.Redirect(res, req, "/welcome", http.StatusFound)
			fmt.Println("1")
		default:
			http.Redirect(res, req, "/dashboard", http.StatusFound)
			fmt.Println("2")
		}
	} else {
		// if val is not a string type
		http.Redirect(res, req, "/welcome", http.StatusFound)
	}
}
