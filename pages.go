package main

import "net/http"
import "io/ioutil"

type Navbar struct {
	Navbar string
}

func DashboardHandler(res http.ResponseWriter, req *http.Request) {
	var ret = CheckIfSessionIsValid(res, req)
	if ret == false {
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
	var ret = CheckIfSessionIsValid(res, req)
	if ret == false {
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
