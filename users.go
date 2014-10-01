package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//"encoding/json"
)

import "io/ioutil"

func UsersHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter WelcomeHandler")

	/*
		var ret = CheckIfSessionIsValid(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/
	switch req.Method {
	case "GET":
		fmt.Println("GET request")
		// Serve the resource.
	case "POST":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Create a new record.
	case "PUT":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Update an existing record.
	case "DELETE":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Remove the record.
	default:
		// Give an error message.
	}

	//data, _ := json.Marshal("{'hello':'wercker!'}")
	content, err := ioutil.ReadFile("db/users.json")
	if err == nil {
		//Do something
		fmt.Println("Cannot open file for reading")
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write([]byte(content))

}

func UserHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter Users")

	vars := mux.Vars(req)
	key := vars["user"]

	fmt.Println("key: %s", key)

	switch req.Method {
	case "GET":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Serve the resource.
	case "POST":
		http.Error(res, "POST not valid", http.StatusInternalServerError)
		return
		// Create a new record.
	case "PUT":
		// Update an existing record.
	case "DELETE":
		// Remove the record.
	default:
		// Give an error message.
	}

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
