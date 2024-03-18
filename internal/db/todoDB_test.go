package db

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	tdDB := CreateTodoDB()
	id := tdDB.createTodo("Hello world", nil, time.Now())

	todo , err := tdDB.getTodo(id)
	if err != nil {
		t.Fatal(err)
	}

	if todo.Description != "Hello world" {
		t.Errorf("got todo.Description=%s, needs %s", todo.Description, "Hello world")
	}

	if todo.Id != id {
		t.Errorf("got todo.Id=%d, id=%d", todo.Id, id)
	}

	allTodo := tdDB.getAllTodos()
	if len(allTodo) != 1 {
		t.Errorf("got len(allTodo) = %d, need 1", len(allTodo))
	}

	_, err = tdDB.getTodo(id + 1)
	if err == nil {
		t.Errorf("got nil, need error")
	}

	tdDB.createTodo("Hey", nil, time.Now())
	allTodo2 := tdDB.getAllTodos()
	if len(allTodo2) != 2 {
		t.Errorf("got len(allTodo2) = %d, need 2", len(allTodo2))
	}
}

func TestDelete(t *testing.T) {
	tdDB := CreateTodoDB()
	id1 := tdDB.createTodo("Wash cloth", nil, time.Now())
	id2 := tdDB.createTodo("Arrange book", nil, time.Now())

	if err := tdDB.deleteTodo(id1 + 100); err == nil {
		t.Fatalf("deleting invalid todo=%d resulting no error, need error", id1 +100)
	}
	if err := tdDB.deleteTodo(id1); err != nil {
		t.Fatal(err)
	}
	if err := tdDB.deleteTodo(id2); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAll(t *testing.T) {
	tdDB := CreateTodoDB()
	tdDB.createTodo("Wash cloth", nil, time.Now())
	tdDB.createTodo("Arrange book", nil, time.Now())

	if err := tdDB.deleteAllTodos(); err != nil {
		t.Fatal(err)
	}

	todo := tdDB.getAllTodos()
	if len(todo) != 0 {
		t.Fatal("unable to delete all todo")
	}
}

func TestGetTodoByTag(t *testing.T) {
	tdDB := CreateTodoDB()
	tdDB.createTodo("Wash plate", []string{"House"}, time.Now())
	tdDB.createTodo("Build go server", []string{"Coding"}, time.Now())
	tdDB.createTodo("Help Mom", []string{"House"}, time.Now())
	tdDB.createTodo("Read Networking book", []string{"Coding"}, time.Now())
	tdDB.createTodo("Setup wifi", []string{"House", "Coding"}, time.Now())

	var tests = []struct {
		tag string
		counts int
	} {
		{"House", 3},
		{"Coding", 3},
		{"Room", 0},
	}
	for _, testCase := range tests {
		t.Run(testCase.tag, func(t *testing.T) {
			numberTagReturned := len(tdDB.GetTodoByTag(testCase.tag))
			 if numberTagReturned != testCase.counts {
				t.Errorf("got %v, need %v for %s", numberTagReturned, testCase.counts, testCase.tag)
			 }
		})
	}
}

func TestGetTasksByDueDate(t *testing.T) {
	timeFormat := "2006-Jan-02"
	parseDate := func(tstr string) time.Time {
		tt, err := time.Parse(timeFormat, tstr)
		if err != nil {
			t.Fatal(err)
		}
		return tt
	}

	tdDB := CreateTodoDB()
	tdDB.createTodo("Wash plate", nil, parseDate("2022-Dec-02"))
	tdDB.createTodo("Build go server", nil, parseDate("2023-Jan-12"))
	tdDB.createTodo("Help Mom", nil, parseDate("2021-Jan-12"))
	tdDB.createTodo("Read Networking book", nil, parseDate("2021-Jan-12"))
	tdDB.createTodo("Setup wifi", nil, parseDate("2021-Jan-12"))

	// retrieve a single todo
	y, m, d := parseDate("2022-Dec-02").Date()
	todoWashPlate := tdDB.getTodosByDueDate(y, m, d)
	if len(todoWashPlate) != 1 {
		t.Errorf("retrieve error %v", todoWashPlate)
	}

	tests := []struct{
		date time.Time
		counts int
	}{
		{parseDate("2022-Dec-02"), 1},
		{parseDate("2023-Jan-12"), 1},
		{parseDate("2021-Jan-12"), 3},
		{parseDate("2024-Mar-14"), 0},
	}

	for _, test := range(tests) {
		t.Run(test.date.String(), func(t *testing.T) {
			todoList := tdDB.getTodosByDueDate(test.date.Date())
			if len(todoList) != test.counts {
				t.Fatalf("Failed at %v, should have %v", test.date.String(), test.counts)
			}
		})
	}
}