package db

import (
	"sync"
	"time"
	"fmt"
)

type Todo struct {
	Id int
	Description string
	Tags []string
	Due time.Time
}

// in-memory DB for Todo
type TodoDB struct {
	m sync.Mutex
	todos map[int]Todo
	nextId int
}

func CreateTodoDB() *TodoDB {
	td := &TodoDB{}
	td.todos = make(map[int]Todo)
	td.nextId = 0
	return td
}

func (td *TodoDB) CreateTodo(description string, tags []string, due time.Time) int {
	td.m.Lock()
	defer td.m.Unlock()

	id := td.nextId
	todoTag := make([]string, len(tags))
	copy(todoTag,tags)
	newTodo := Todo{id, description, todoTag, due}

	td.todos[id] = newTodo
	td.nextId++

	return id
}

func (td *TodoDB) GetTodo(id int) (Todo, error) {
	td.m.Lock()
	defer td.m.Unlock()

	todo, ok := td.todos[id]
	if !ok {
		return Todo{}, fmt.Errorf("invalid id cannot find id=%d", id)
	} else {
		return todo, nil
	}
}

func (tdDB *TodoDB) DeleteTodo(id int) error {
	tdDB.m.Lock()
	defer tdDB.m.Unlock()

	if _, ok := tdDB.todos[id]; !ok {
		return fmt.Errorf("todo not found id = %d", id)
	}
	delete(tdDB.todos, id)
	return nil
}

func (tdDB *TodoDB) DeleteAllTodos() error {
	tdDB.m.Lock()
	defer tdDB.m.Unlock()

	tdDB.todos = make(map[int]Todo)
	return nil
}

func (tdDB *TodoDB) GetAllTodos() []Todo {
	tdDB.m.Lock()
	defer tdDB.m.Unlock()

	todoList := make([]Todo, 0, len(tdDB.todos))

	for _, todo := range(tdDB.todos) {
		todoList = append(todoList, todo)
	}

	return todoList
}

func (tdDB *TodoDB) GetTodoByTag(tag string) []Todo {
	tdDB.m.Lock()
	defer tdDB.m.Unlock()

	var todoList []Todo
	todoloop:
		for _, todo := range tdDB.todos {
			for _, todoTag := range todo.Tags {
				if todoTag == tag {
					todoList = append(todoList, todo)
					continue todoloop
				}
			}
		}
	return todoList
}

func (tdDB *TodoDB) GetTodosByDueDate(year int, month time.Month, day int) []Todo {
	tdDB.m.Lock()
	defer tdDB.m.Unlock()

	var todoList []Todo

	for _, todo := range(tdDB.todos) {
		y, m, d := todo.Due.Date()

		if y == year && m == month && d == day {
			todoList = append(todoList, todo)
		}
	}

	return todoList
}	
