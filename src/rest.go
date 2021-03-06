package main

import "os"
import "fmt"
import "github.com/gorilla/mux"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "errors"
import "time"
import "bytes"
import "bufio"

func userHanderGet(ctx *SessionCtx, res http.ResponseWriter, req *http.Request) {
    path := fmt.Sprintf("%s/%s", ctx.Db, "users.json")
	content, err := ioutil.ReadFile(path)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		fmt.Println("Cannot open file for reading %s", err)
		return
	}
	res.Write([]byte(content))
}

func RestUsersHandler(res http.ResponseWriter, req *http.Request) {
	ctx, ok := getSessionCtx(req)
    if !ok {
		http.Error(res, "Not authorized, sorry", http.StatusUnauthorized)
        return
    }

	switch req.Method {
	case "GET":
		userHanderGet(&ctx, res, req)
		return
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
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
		// Give an error message.
	}

	http.Error(res, "internal error", http.StatusInternalServerError)
	return
}

func RestUserHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Println("enter RestUserHandler")

	vars := mux.Vars(req)
	key := vars["user"]

	fmt.Println("key: ", key)

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
		var ret = getSessionCtxalidSession(res, req)
		if ret == false && req.URL.Path != "/signInP" {
				http.Redirect(res, req, "/signInP", http.StatusFound)
				return
		}
	*/

	p, err := loadPageTemplate("welcome")
	if err != nil {
		http.Error(res, "foooooo", http.StatusInternalServerError)
	}
	fmt.Println(p)
	x := Person{Name: "Mary"}
	p.Execute(res, x)
}

const itemsFilePath string = "items.json"

func userItemFile(ctx *SessionCtx) (string, bool) {
    if ctx.Db == "" {
        panic("no path given")
    }
    return fmt.Sprintf("%s/%s", ctx.Db, itemsFilePath), true
}

func generateAllItemsAsJson(ctx *SessionCtx, r *http.Request) (data string, err error) {
    filePath, ok := userItemFile(ctx)
    if ok != true {
        panic("need fix")
    }

    fmt.Println("open item file at", filePath)
    file, err := os.Open(filePath)
	if err != nil {
        fmt.Println("No item data file available yet (no data added)")
		return "", err
	}
	defer file.Close()

	var buffer bytes.Buffer
	var is_not_first bool = false
	buffer.WriteString("[\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if is_not_first {
			buffer.WriteString(",")
		}
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
		is_not_first = true
	}
	buffer.WriteString("]\n")
	return buffer.String(), scanner.Err()
}

func itemsHanderGet(ctx *SessionCtx, res http.ResponseWriter, req *http.Request) {
	data, err := generateAllItemsAsJson(ctx, req)
	if err != nil {
		// if an error occur we return an empty JSON array
		data = "[ ]"
	}
	fmt.Println(data)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write([]byte(data))
}

// JSON encoded content in form of
// { desc: "foo", bar: "baz" }
// no comas are added, go logic will
// add these as a prepartion before
// transport
type item_file_line struct {
	Id               string
	Description      string
	Deadline         string
	Priority         string
	AssignedTo       string
	AssociatedPerson []string
	Tags             []string
	CreationDate     string
	ModifiedDate     string
	Information      string
}

type item_struct struct {
	Command string
	Data    map[string]string
}

type ItemJson struct {
	Command      string       `json:Command"`
	ItemJsonData ItemJsonData `json:"Data"`
}

type ItemJsonData struct {
	Description string   `json:"Description"`
	Deadline    string   `json:"Deadline"`
	AssignedTo  string   `json:"AssignedTo"`
	Priority    string   `json:"Priority"`
	Information string   `json:"Information"`
	Tags        []string `json:"Tags"`
	Persons     []string `json:"AssociatedPersons"`
}

//func checkIfDataisValid(data item_struct) error {
// read over all items and get the highest
// item id, we will get the new item id+1
// item recycling is not implemented now
//	desc, ok := data.Data["Description"]
//	if ok != true {
//		return errors.New("Item Description mission from struct")
//	}

//	fmt.Println(desc)

//	return nil
//}

// Open file in append mode and add
// JSON encoded line
func appendItemData(ctx *SessionCtx, data item_file_line) error {
    filepath, _ := userItemFile(ctx)
    _, err := os.Stat(filepath)
	if err != nil {
		// no such file or dir
	}

    f, err := os.OpenFile(filepath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s, _ := json.Marshal(data)
	if _, err = f.WriteString(string(s) + "\n"); err != nil {
		panic(err)
	}

	return nil
}

// return UNIX time in milliseconds,
// this should be unique enough for one user
// for now
func getNewItemId() string {
	ct := time.Now()
	now := ct.Nanosecond()
	miliSeconds := (now % 1e9) / 1e6
	sec := ct.UTC().Format("20060102150405")
	return fmt.Sprintf("item-%s%03d", sec, miliSeconds)
}

func addItem(ctx *SessionCtx, data ItemJson) error {

	var err error

	if data.Command != "add" {
		// the only command I know, return if not add
		return errors.New("only add supported")
	}

	//	err := checkIfDataisValid(data)
	//	if err != nil {
	//		return err
	//	}

	var new_data item_file_line
	new_data.Id = getNewItemId()
	new_data.Description = data.ItemJsonData.Description
	new_data.Deadline = data.ItemJsonData.Deadline
	new_data.Priority = data.ItemJsonData.Priority
	new_data.AssignedTo = data.ItemJsonData.AssignedTo
	new_data.AssociatedPerson = data.ItemJsonData.Persons
	new_data.Tags = data.ItemJsonData.Tags
	new_data.Information = data.ItemJsonData.Information

	new_data.CreationDate = time.Now().UTC().Format("20060102150405")
	new_data.ModifiedDate = new_data.CreationDate

	err = appendItemData(ctx, new_data)
	if err != nil {
		return err
	}

	return nil
}

type client_items_response_msg struct {
	Status string
}

func itemsHanderPost(ctx *SessionCtx, w http.ResponseWriter, r *http.Request) {
	var t ItemJson

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failure in reading from client %s", err)
		// FIXME: return negative status code
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("Failure in item POST struct", err)
		return
	}

	fmt.Println(t.Command)
	err = errors.New("Not implemented")
	switch t.Command {
	case "add":
		err = addItem(ctx, t)
		// Serve the resource.
	case "del":
	default:
		// Give an error message.
	}

	if err != nil {
		fmt.Println("item error occured", err)
		return
	}

	var res_msg client_items_response_msg
	res_msg.Status = "success"

	msg, _ := json.Marshal(res_msg)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(msg))
}

func RestItemsHandler(res http.ResponseWriter, req *http.Request) {
    sessionCtx, ok := getSessionCtx(req)
	if ok == false {
        panic("not authenfificated")
		return
	}

	switch req.Method {
	case "GET":
		itemsHanderGet(&sessionCtx, res, req)
		return
	case "POST":
		itemsHanderPost(&sessionCtx, res, req)
		return
	case "PUT":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
	case "DELETE":
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
	default:
		http.Error(res, "Not allows", http.StatusInternalServerError)
		return
	}

	http.Error(res, "internal error", http.StatusInternalServerError)
	return
}

/*

project JSON data

type: [ project | task | deadline ]
id: [ project-ID | task-ID | deadline-ID

project has no start or endtime. THe starttime is the min
of all tasks/deadlines. But is never displayed somewhere.
There can be several projects

tasks have a start and enddata, there must be at least a differnce
of one day

deadlines have a start and enddata with exactly the same data

tasks can be assigned to exactly one project

deadlines can be assigned to exactly one project

tasks can have one or multiple users

projects, tasks, and deadlines must have a description

*/

type ProjectGetCmdJson struct {
	Command      string       `json:command"`
}

func projectsHanderPost(ctx *SessionCtx, w http.ResponseWriter, r *http.Request) {
	var t ProjectGetCmdJson
	fmt.Println("xxx")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failure in reading from client %s", err)
		return
	}


	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println("Failure in item POST struct", err)
		return
	}

	fmt.Println(t.Command)
	err = errors.New("Not implemented")
	switch t.Command {
	case "all":
		fmt.Println("all")
		return
	default:
	}

	panic("")
}

func RestProjectsHandler(res http.ResponseWriter, req *http.Request) {
    sessionCtx, ok := getSessionCtx(req)
	if ok == false {
        panic("not authenfificated")
		return
	}

	fmt.Println("Project handler")

	switch req.Method {
	case "GET":
		http.Error(res, "Not allowed", http.StatusInternalServerError)
		return
	case "POST":
		projectsHanderPost(&sessionCtx, res, req)
		return
	case "PUT":
		http.Error(res, "Not allowed", http.StatusInternalServerError)
		return
	case "DELETE":
		http.Error(res, "Not allowed", http.StatusInternalServerError)
		return
	default:
		http.Error(res, "Not allowed", http.StatusInternalServerError)
		return
	}

	http.Error(res, "Iinternal error", http.StatusInternalServerError)
	return
}
