package service

import (
	"github.com/yesseneon/todo"
	"github.com/yesseneon/todo/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userID, listID int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetByID(userID, listID)
	if err != nil {
		// list doesn't exists or doesn't belongs to user
		return 0, err
	}

	return s.repo.Create(listID, item)
}

func (s *TodoItemService) GetAll(userID, listID int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userID, listID)
}

func (s *TodoItemService) GetByID(userID, itemID int) (todo.TodoItem, error) {
	return s.repo.GetByID(userID, itemID)
}

func (s *TodoItemService) Update(userID, itemID int, item todo.TodoItemUpdate) error {
	if err := item.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userID, itemID, item)
}

func (s *TodoItemService) Delete(userID, itemID int) error {
	return s.repo.Delete(userID, itemID)
}
