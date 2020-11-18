package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yesseneon/todo"
)

type TodoItemRepository struct {
	db *sqlx.DB
}

func newTodoItemRepository(db *sqlx.DB) *TodoItemRepository {
	return &TodoItemRepository{db: db}
}

func (r *TodoItemRepository) Create(listID int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListsItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListsItemQuery, listID, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoItemRepository) GetAll(userID, listID int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.* FROM %s ti
		INNER JOIN %s li ON ti.id=li.item_id
		INNER JOIN %s ul ON ul.list_id=li.list_id
		WHERE li.list_id=$1 AND ul.user_id=$2`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Select(&items, query, listID, userID)
	return items, err
}

func (r *TodoItemRepository) GetByID(userID, itemID int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.* FROM %s ti
		INNER JOIN %s li ON ti.id=li.item_id
		INNER JOIN %s ul ON ul.list_id=li.list_id
		WHERE ti.id=$1 AND ul.user_id=$2`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, itemID, userID)
	return item, err
}

func (r *TodoItemRepository) Update(userID, itemID int, item todo.TodoItemUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *item.Title)
		argID++
	}

	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *item.Description)
		argID++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argID))
		args = append(args, *item.Done)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
		WHERE ti.id=li.item_id AND li.list_id=ul.list_id
		AND ul.user_id=$%d AND ti.id=$%d`, todoItemsTable, setQuery, listsItemsTable, usersListsTable, argID, argID+1)
	args = append(args, userID, itemID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemRepository) Delete(userID, itemID int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
		WHERE ti.id=li.item_id AND li.list_id=ul.list_id
		AND ul.user_id=$1 AND ti.id=$2`, todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, itemID)
	return err
}
