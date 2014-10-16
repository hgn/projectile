package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)
import "encoding/json"

type SessionCtx struct {
	Username string
	Db       string
}

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
	router.HandleFunc("/signup", SignUpHandler)

	router.HandleFunc("/dashboard", DashboardHandler)
	router.HandleFunc("/items", ItemsHandler)

	router.HandleFunc("/show", ShowHandler)
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

func getSessionCtx(req *http.Request) (SessionCtx, bool) {
	var sessionCtx SessionCtx
	session, _ := store.Get(req, sessionName)
	if val, ok := session.Values["username"].(string); ok {
		// if val is a string
		switch val {
		case "":
			return sessionCtx, false
		default:
			sessionCtx.Username = val
			sessionCtx.Db, _ = session.Values["Db"].(string)
			return sessionCtx, true
		}
	} else {
		return sessionCtx, false
	}
}

var dbPasswdData []dbPasswdDataEntry

type dbPasswdDataEntry struct {
	Username string `json:"Username"`
	Photo    string `json:"Photo"`
	Db       string `json:"Db"`
	Password string `json:"Password"`
}

func userValid(username, password string) (dbPasswdDataEntry, bool) {
	var rentry dbPasswdDataEntry
	configFile, err := os.Open("db/passwd.json")
	if err != nil {
		fmt.Println("opening config file failed", err.Error())
		return rentry, false
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&dbPasswdData); err != nil {
		fmt.Println("parsing config file", err.Error())
		return rentry, false
	}

	// search hashed password
	hashedPassword := ""
	for _, entry := range dbPasswdData {
		if entry.Username == username {
			rentry = entry
			hashedPassword = entry.Password
			break
		}
	}
	if hashedPassword == "" {
		fmt.Println("User not in database or password not available")
		newPassword, err := CryptPassword([]byte(username))
		if err == nil {
			// we do not handle errors here
			// because this is just a workaround
			// functionality for some time
			fmt.Printf("Possible new hashed passwort for %s: %s\n",
				username, newPassword)
		}
		return rentry, false
	}

	return rentry, CheckPassword([]byte(password), []byte(hashedPassword))
}

const dbRootPath string = "db"

func dbPath(entry *dbPasswdDataEntry) string {
    return fmt.Sprintf("%s/%s-%s", dbRootPath, entry.Username, entry.Db)
}

func checkUserEnvironment(entry *dbPasswdDataEntry) {
    dbpath := dbPath(entry)

   src, err := os.Stat(dbpath)
   if err != nil {
       err = os.MkdirAll(dbpath, 0777)
       if err != nil {
           fmt.Println("Something failed creating directory", dbpath)
       }
   }

   if src.IsDir() {
     fmt.Println("User specific directory already created, nothing to do")
     return
   }

   fmt.Println("Userdirectory not available, create it now", dbpath)
   err = os.MkdirAll(dbpath, 0777)
   if err != nil {
       fmt.Println("Something failed creating directory", dbpath)
   }
}

//handler for signIn
func SignInHandler(res http.ResponseWriter, req *http.Request) {
	username, password := req.FormValue("username"), req.FormValue("password")

	dbPasswdDataEntry, ret := userValid(username, password)
	if ret == false {
		fmt.Printf("Password for user %s invalid\n", username)
		http.Redirect(res, req, "/welcome", http.StatusFound)
		return
	}

	sessionNew, _ := store.Get(req, sessionName)

	sessionNew.Values["username"] = username
	sessionNew.Values["Photo"] = dbPasswdDataEntry.Photo
	sessionNew.Values["Db"] = dbPath(&dbPasswdDataEntry)

    checkUserEnvironment(&dbPasswdDataEntry)

	err := sessionNew.Save(req, res)
	if err != nil {
		panic("session save error")
	}
	store.Save(req, res, sessionNew)

	http.Redirect(res, req, "/dashboard", http.StatusFound)
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

func loadPage(title string) ([]byte, error) {
	filename := "page-templates/" + title + ".html"
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return dat, nil
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
