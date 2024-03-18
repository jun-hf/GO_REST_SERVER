package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"
	"github.com/jun-hf/GO_REST_SERVER/internal/db"
)

type ToDoServer struct {
	todoDB *db.TodoDB
}

func createTodoServer() *ToDoServer {
	tdDB := db.CreateTodoDB()
	return &ToDoServer{todoDB: tdDB}
}

func (tdServer *ToDoServer) createTodoHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("processing todo create at %s\n", req.URL.Path)

	type RequestTodo struct {
		Description string
		Tags        []string
		Due         time.Time
	}

	type ResponseId struct {
		Id int
	}

	// ensure is application/json
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "expecting application/json", http.StatusUnsupportedMediaType)
		return
	}

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	var reqTodo RequestTodo
	if err := decode.Decode(&reqTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := tdServer.todoDB.CreateTodo(reqTodo.Description, reqTodo.Tags, reqTodo.Due)
	
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (tdServer *ToDoServer) getAllTodosHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("processing getAllTodosHandler at %s\n", req.URL.Path)

	todosList := tdServer.todoDB.GetAllTodos()
	js, err := json.Marshal(todosList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (tdServer *ToDoServer) robotHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("processing robotHandler at %s\n", req.URL.Path)
	w.Write([]byte("Heelo"))
}

func (tdServer *ToDoServer) getTodoHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling getTodoHandler at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := tdServer.todoDB.GetTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (tdServer *ToDoServer) deleteTodoHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling deleteTodoHandler at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id please provide a valid integer", http.StatusBadRequest)
		return
	}

	err = tdServer.todoDB.DeleteTodo(id)
	if err != nil {
		http.Error(w, "id does not exist in db", http.StatusBadRequest)
	}
}

func (tdServer *ToDoServer) deleteAllTodosHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Deleting all todos at %s\n", req.URL.Path)
	tdServer.todoDB.DeleteAllTodos()
}

func (tdServer *ToDoServer) getByTagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Getting by tags at %s\n", req.URL.Path)
	tag := req.PathValue("tag")

	todos := tdServer.todoDB.GetTodoByTag(tag)
	
	js, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (tdServer *ToDoServer) getByDueDateHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Getting By Due date %s\n", req.URL.Path)

	writeJsonResponse := func(w http.ResponseWriter, response []byte) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))
	if errYear != nil || errMonth !=nil || errDay !=nil || month < int(time.January) || month > int(time.December) {
		http.Error(w, "Invalid request format expect: /year/month/day in int", http.StatusBadRequest)
		return
	}

	todos := tdServer.todoDB.GetTodosByDueDate(year, time.Month(month), day)
	js, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, js)
}

func main() {
	mux := http.NewServeMux()
	todoServer := createTodoServer()

	mux.HandleFunc("POST /todo/", todoServer.createTodoHandler)
	mux.HandleFunc("GET /todos/", todoServer.getAllTodosHandler)
	mux.HandleFunc("GET /robot.txt", todoServer.robotHandler)
	mux.HandleFunc("GET /todo/{id}", todoServer.getTodoHandler)
	mux.HandleFunc("DELETE /todo/{id}", todoServer.deleteTodoHandler)
	mux.HandleFunc("DELETE /deleteAllTodos", todoServer.deleteAllTodosHandler)
	mux.HandleFunc("GET /tag/{tag}", todoServer.getByTagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}", todoServer.getByDueDateHandler)

	http.ListenAndServe(":3000", mux)
}
