package main

import "fmt"
import "net/http"
import "io/ioutil"


type Navbar struct {
	Navbar string
}


func DashboardHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter DashboardHandler")

	var ret = CheckIfSessionIsValid(res, req)
	if ret == false {
		http.Redirect(res, req, "/signInP", http.StatusFound)
		return
	}

	p, err := loadPage("dashboard")
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

	fmt.Println("enter ItemsHandler")

	var ret = CheckIfSessionIsValid(res, req)
	if ret == false {
		http.Redirect(res, req, "/signInP", http.StatusFound)
		return
	}

	p, err := loadPage("items")
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
