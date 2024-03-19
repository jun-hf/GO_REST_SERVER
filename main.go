package main

import (
	"net/http"
	"strconv"
	"time"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/jun-hf/GO_REST_SERVER/internal/db"
)

type ToDoServer struct {
	todoDB *db.TodoDB
}

func createTodoServer() *ToDoServer {
	tdDB := db.CreateTodoDB()
	return &ToDoServer{todoDB: tdDB}
}

func (tdServer *ToDoServer) getAllTodoHandler(c *gin.Context) {
	allTodos := tdServer.todoDB.GetAllTodos()
	c.JSON(http.StatusOK, allTodos)
}

func (tdServer *ToDoServer) deleteAllTodosHandler(c *gin.Context) {
	tdServer.todoDB.DeleteAllTodos()
}

func (tdServer *ToDoServer) createTodoHandler(c *gin.Context) {
	type RequestTodo struct {
		Description string `json:"description"`
		Tags []string `json:"tags"`
		Due time.Time `json:"due"`
	}

	var reqTodo RequestTodo
	if err := c.ShouldBindJSON(&reqTodo); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := tdServer.todoDB.CreateTodo(reqTodo.Description, reqTodo.Tags, reqTodo.Due)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (tdServer *ToDoServer) getTodoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	todo, err := tdServer.todoDB.GetTodo(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (tdServer *ToDoServer) deleteTodoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = tdServer.todoDB.DeleteTodo(id); err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
}

func (tdServer *ToDoServer) tagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")
	todos := tdServer.todoDB.GetTodoByTag(tag)
	c.JSON(http.StatusOK, todos)
}

func (tdServer *ToDoServer) dueHandler(c *gin.Context) {
	year, errYear := strconv.Atoi(c.Params.ByName("year"))
	month, errMonth := strconv.Atoi(c.Params.ByName("month"))
	day, errDay := strconv.Atoi(c.Params.ByName("day"))
	if errYear != nil || errMonth !=nil || errDay !=nil || month < int(time.January) || month > int(time.December) {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>, got %v", c.FullPath())
		return
	}

	todos := tdServer.todoDB.GetTodosByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, todos)
}

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	server := createTodoServer()

	router.POST("/todos", server.createTodoHandler)
	router.GET("/todos", server.getAllTodoHandler)
	router.DELETE("/todos/deleteAll", server.deleteAllTodosHandler)
	router.GET("/todos/{id}", server.getTodoHandler)
	router.DELETE("/todos/{id}", server.deleteTodoHandler)
	router.GET("/tag", server.tagHandler)
	router.GET("/due/{year}/{month}/{day}", server.dueHandler)

	router.Run("localhost:" + os.Getenv("PORT"))
}
