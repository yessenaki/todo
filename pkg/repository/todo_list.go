package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yesseneon/todo"
)

type TodoListRepository struct {
	db *sqlx.DB
}

func newTodoListRepository(db *sqlx.DB) *TodoListRepository {
	return &TodoListRepository{db: db}
}

func (r *TodoListRepository) Create(userID int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userID, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListRepository) GetAll(userID int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("SELECT tl.* FROM %s tl INNER JOIN %s ul ON tl.id=ul.list_id WHERE ul.user_id=$1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userID)
	return lists, err
}

func (r *TodoListRepository) GetByID(userID, listID int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf("SELECT tl.* FROM %s tl INNER JOIN %s ul ON tl.id=ul.list_id WHERE ul.user_id=$1 AND ul.list_id=$2", todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userID, listID)
	return list, err
}

func (r *TodoListRepository) Update(userID, listID int, list todo.TodoListUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if list.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *list.Title)
		argID++
	}

	if list.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *list.Description)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id=ul.list_id AND ul.user_id=$%d AND ul.list_id=$%d", todoListsTable, setQuery, usersListsTable, argID, argID+1)
	args = append(args, userID, listID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoListRepository) Delete(userID, listID int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id=ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, listID)
	return err
}
