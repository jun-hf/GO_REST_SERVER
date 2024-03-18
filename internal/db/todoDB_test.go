package db

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	tdDB := CreateTodoDB()
	id := tdDB.CreateTodo("Hello world", nil, time.Now())

	todo , err := tdDB.GetTodo(id)
	if err != nil {
		t.Fatal(err)
	}

	if todo.Description != "Hello world" {
		t.Errorf("got todo.Description=%s, needs %s", todo.Description, "Hello world")
	}

	if todo.Id != id {
		t.Errorf("got todo.Id=%d, id=%d", todo.Id, id)
	}

	allTodo := tdDB.GetAllTodos()
	if len(allTodo) != 1 {
		t.Errorf("got len(allTodo) = %d, need 1", len(allTodo))
	}

	_, err = tdDB.GetTodo(id + 1)
	if err == nil {
		t.Errorf("got nil, need error")
	}

	tdDB.CreateTodo("Hey", nil, time.Now())
	allTodo2 := tdDB.GetAllTodos()
	if len(allTodo2) != 2 {
		t.Errorf("got len(allTodo2) = %d, need 2", len(allTodo2))
	}
}

func TestDelete(t *testing.T) {
	tdDB := CreateTodoDB()
	id1 := tdDB.CreateTodo("Wash cloth", nil, time.Now())
	id2 := tdDB.CreateTodo("Arrange book", nil, time.Now())

	if err := tdDB.DeleteTodo(id1 + 100); err == nil {
		t.Fatalf("deleting invalid todo=%d resulting no error, need error", id1 +100)
	}
	if err := tdDB.DeleteTodo(id1); err != nil {
		t.Fatal(err)
	}
	if err := tdDB.DeleteTodo(id2); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAll(t *testing.T) {
	tdDB := CreateTodoDB()
	tdDB.CreateTodo("Wash cloth", nil, time.Now())
	tdDB.CreateTodo("Arrange book", nil, time.Now())

	if err := tdDB.DeleteAllTodos(); err != nil {
		t.Fatal(err)
	}

	todo := tdDB.GetAllTodos()
	if len(todo) != 0 {
		t.Fatal("unable to delete all todo")
	}
}

func TestGetTodoByTag(t *testing.T) {
	tdDB := CreateTodoDB()
	tdDB.CreateTodo("Wash plate", []string{"House"}, time.Now())
	tdDB.CreateTodo("Build go server", []string{"Coding"}, time.Now())
	tdDB.CreateTodo("Help Mom", []string{"House"}, time.Now())
	tdDB.CreateTodo("Read Networking book", []string{"Coding"}, time.Now())
	tdDB.CreateTodo("Setup wifi", []string{"House", "Coding"}, time.Now())

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
				t.Error(tdDB.todos)
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
	tdDB.CreateTodo("Wash plate", nil, parseDate("2022-Dec-02"))
	tdDB.CreateTodo("Build go server", nil, parseDate("2023-Jan-12"))
	tdDB.CreateTodo("Help Mom", nil, parseDate("2021-Jan-12"))
	tdDB.CreateTodo("Read Networking book", nil, parseDate("2021-Jan-12"))
	tdDB.CreateTodo("Setup wifi", nil, parseDate("2021-Jan-12"))

	// retrieve a single todo
	y, m, d := parseDate("2022-Dec-02").Date()
	todoWashPlate := tdDB.GetTodosByDueDate(y, m, d)
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
			todoList := tdDB.GetTodosByDueDate(test.date.Date())
			if len(todoList) != test.counts {
				t.Fatalf("Failed at %v, should have %v", test.date.String(), test.counts)
			}
		})
	}
}