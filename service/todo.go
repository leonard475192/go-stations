package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	var rows *sql.Rows
	var err error
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
		if err != nil {
			// log.Printf("Error QueryContext not prevID:%v", err)
			return nil, err
		}
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			// log.Printf("Error QueryContext:%v", err)
			return nil, err
		}
	}

	TODOs := make([]*model.TODO, 0)
	for rows.Next() {
		var TODO model.TODO
		err := rows.Scan(
			&TODO.ID,
			&TODO.Subject,
			&TODO.Description,
			&TODO.CreatedAt,
			&TODO.UpdatedAt,
		)
		if err != nil {
			// log.Printf("Error Scan:%v", err)
			return nil, err
		}
		TODOs = append(TODOs, &TODO)
	}

	return TODOs, nil
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
	if len(ids) == 0 {
		return nil
	}

	delete := fmt.Sprintf(deleteFmt, strings.Repeat(", ?", len(ids)-1))
	prepare_todos, err := s.db.PrepareContext(ctx, delete)
	if err != nil {
		log.Printf("Error PrepareContext:%v", err)
		return err
	}
	var list = make([]interface{}, 0)
	for _, id := range ids {
		list = append(list, id)
	}
	delete_todos, err := prepare_todos.ExecContext(ctx, list...)
	if err != nil {
		// log.Printf("Error delete:%v", err)
		return err
	}
	delete_rows, err := delete_todos.RowsAffected()
	if err != nil {
		// log.Printf("Error delete:%v", err)
		return err
	}
	if delete_rows == 0 {
		return model.ErrNotFound{}
	}

	return nil
}
