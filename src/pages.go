package main

import "net/http"
import "io/ioutil"
import "fmt"

type Navbar struct {
	Navbar string
}

func SignUpHandler(res http.ResponseWriter, req *http.Request) {
	p, err := loadPage("signup")
	if err != nil {
		http.Error(res, "Failed to load signin page", http.StatusInternalServerError)
		fmt.Println("Failed to load page from file: %v", err)
	}
	res.Write(p)
}

type Context struct {
	Username  string
	PhotoPath string
	DbPath    string
}

func DashboardHandler(res http.ResponseWriter, req *http.Request) {
    _, ok := getSessionCtx(req)
	if ok == false {
		http.Redirect(res, req, "/signInP", http.StatusFound)
		return
	}

	p, err := loadPageTemplate("dashboard")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}

	content, err := ioutil.ReadFile("page-templates/navbar.html")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}

	x := Navbar{Navbar: string(content[:])}
	p.Execute(res, x)
}

func ItemsHandler(res http.ResponseWriter, req *http.Request) {
	_, ok := getSessionCtx(req)
	if ok == false {
		http.Redirect(res, req, "/signInP", http.StatusFound)
		return
	}

	p, err := loadPageTemplate("items")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}

	content, err := ioutil.ReadFile("page-templates/navbar.html")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}

	x := Navbar{Navbar: string(content[:])}
	p.Execute(res, x)
}
