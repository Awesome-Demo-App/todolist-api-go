package service

import (
	"errors"

	domain "github.com/awesome-demo-app/todolist-api/core/domain"
	ports "github.com/awesome-demo-app/todolist-api/core/ports"
)

type todoService struct {
	todoRepository ports.ToDoRepository
}

func New(repository ports.ToDoRepository) *todoService {
	return &todoService{
		todoRepository: repository,
	}
}

func (service *todoService) GetAll() ([]domain.ToDo, error) {
	todos, err := service.todoRepository.GetAll()

	if err != nil {
		return nil, errors.New("getting todos from repository has failed")
	}

	return todos, nil
}

func (service *todoService) Create(summary string) domain.ToDo {
	fmt.Printf("Domain → Creating a new ToDo with summary: %s\n", summary)
	return service.todoRepository.Create(domain.ToDo{
		Summary:   summary,
		Completed: false,
	})
}

func (service *todoService) Delete(id uint) {
	fmt.Printf("Domain → Deleting ToDo %d\n", id)
	service.todoRepository.Delete(id)
}
