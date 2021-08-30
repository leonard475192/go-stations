package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/leonard475192/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result_insert, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	new_todo_id, err := result_insert.LastInsertId()
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	result_confirms, err := s.db.QueryContext(ctx, confirm, new_todo_id)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	var new_todo model.TODO
	for result_confirms.Next() {
		if err := result_confirms.Scan(&new_todo.Subject, &new_todo.Description, &new_todo.CreatedAt, &new_todo.UpdatedAt); err != nil {
			log.Fatal(err)
		}
	}

	return &new_todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	return nil, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
