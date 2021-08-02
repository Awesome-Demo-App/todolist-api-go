package ports

import (
	domain "github.com/awesome-demo-app/todolist-api/core/domain"
)

type ToDoRepository interface {
	GetAll() ([]domain.ToDo, error)
	Create(domain.ToDo) domain.ToDo
	Delete(uint)
}

type ToDoService interface {
	GetAll() ([]domain.ToDo, error)
	Create(string) domain.ToDo
	Delete(uint)
}
