package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
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

	id := tdServer.todoDB.createTodo(reqTodo.Description, reqTodo.Tags, reqTodo.Due)
	
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	mux := http.NewServeMux()
	todoServer := createTodoServer()

	mux.HandleFunc("POST /task/", todoServer.createTodoHandler)
}
