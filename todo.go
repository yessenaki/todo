package todo

import "errors"

type TodoList struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	ID     int
	UserID int
	ListID int
}

type TodoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListsItem struct {
	ID     int
	ListID int
	ItemID int
}

type TodoListUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (list TodoListUpdate) Validate() error {
	if list.Title == nil && list.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
