package task

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/pkg/param"
)

var (
	createTaskQuery = `INSERT INTO tasks
		(title, description, status, user_id)
		VALUES ($1, $2, $3, $4) RETURNING *`

	getTaskByUserIDQuery = `SELECT 
		id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY id LIMIT $2 OFFSET $3`

	getTaskByIDQuery = `SELECT 
		id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2`

	updateTaskQuery = `UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6`

	deleteTaskQuery = `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
)

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
	GetByUserID(ctx context.Context, userId int64, param param.Param) ([]model.Task, error)
	GetByID(ctx context.Context, id, userId int64) (model.Task, error)
	Update(ctx context.Context, task model.Task, userId int64) (model.Task, error)
	Delete(ctx context.Context, id, userId int64) error
	Count(ctx context.Context) (int64, error)
}

type Task struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) TaskRepository {
	return &Task{
		db: db,
	}
}

func (t *Task) Create(ctx context.Context, task model.Task) (model.Task, error) {
	result := model.Task{}
	err := t.db.Get(&result, createTaskQuery, task.Title, task.Description, task.Status, task.UserID)
	if err != nil {
		return model.Task{}, err
	}

	return result, nil
}

func (t *Task) GetByUserID(ctx context.Context, userId int64, param param.Param) ([]model.Task, error) {
	result := []model.Task{}

	err := t.db.Select(&result, getTaskByUserIDQuery, userId, param.Limit, param.CalculateOffset())
	if err != nil {
		return []model.Task{}, err
	}

	return result, nil
}

func (t *Task) GetByID(ctx context.Context, id, userId int64) (model.Task, error) {
	result := model.Task{}

	err := t.db.Get(&result, getTaskByIDQuery, id, userId)
	if err != nil {
		return model.Task{}, err
	}

	return result, nil
}

func (t *Task) Update(ctx context.Context, task model.Task, userId int64) (model.Task, error) {
	result, err := t.db.Exec(updateTaskQuery, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID, userId)
	if err != nil {
		return model.Task{}, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (t *Task) Delete(ctx context.Context, id, userId int64) error {
	result, err := t.db.Exec(deleteTaskQuery, id, userId)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Count(ctx context.Context) (int64, error) {
	var total int64
	err := t.db.Get(&total, `SELECT count(*) FROM tasks;`)
	return total, err
}
