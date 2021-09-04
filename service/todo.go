package service

import (
	"context"
	"database/sql"

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

	// TODO ここ聞く
	result_insert, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		// log.Printf("Error ExecContext:%v", err)
		return nil, err
	}
	new_todo_id, err := result_insert.LastInsertId()
	if err != nil {
		// log.Printf("Error LastInsertId:%v", err)
		return nil, err
	}

	confirm_todo := model.TODO{
		ID: new_todo_id,
	}
	err = s.db.QueryRowContext(ctx, confirm, new_todo_id).Scan(
		&confirm_todo.Subject,
		&confirm_todo.Description,
		&confirm_todo.CreatedAt,
		&confirm_todo.UpdatedAt,
	)
	if err != nil {
		// log.Printf("Error QueryContext:%v", err)
		return nil, err
	}

	return &confirm_todo, nil
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
	// ここから
	// TODO を変更する際に利用するメソッドは、PrepareContext メソッドや ExecContext メソッドになり、
	// 保存するTODOを読み取る際に利用するメソッドは QueryRowContext メソッドを利用すると実装することができます。
	prepare_todo, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		// log.Printf("Error PrepareContext:%v", err)
		return nil, model.ErrNotFound{}
	}
	_, err = prepare_todo.ExecContext(ctx, subject, description, id)
	if err != nil {
		// log.Printf("Error Update:%v", err)
		return nil, err
	}
	confirm_todo := model.TODO{
		ID: id,
	}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(
		&confirm_todo.Subject,
		&confirm_todo.Description,
		&confirm_todo.CreatedAt,
		&confirm_todo.UpdatedAt,
	)
	if err != nil {
		// log.Printf("Error Confirm:%v", err)
		return nil, err
	}

	return &confirm_todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
