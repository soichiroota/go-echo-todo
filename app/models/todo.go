package models

import (
	"sample/db"

	"context"
	"time"

	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Content   string    `bun:"content,notnull"`
	Done      bool      `bun:"done"`
	Until     time.Time `bun:"until,nullzero"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

func NewCreateTodoTable() error {
	ctx := context.Background()
	_, err := db.DB.NewCreateTable().Model((*Todo)(nil)).IfNotExists().Exec(ctx)
	return err
}

func GetAllTodos() ([]Todo, error) {
	var todos []Todo
	ctx := context.Background()
	err := db.DB.NewSelect().Model(&todos).Order("created_at").Scan(ctx)
	return todos, err
}

func CreateTodo(todo Todo) (*Todo, error) {
	ctx := context.Background()
	_, err := db.DB.NewInsert().Model(&todo).Exec(ctx)
	return &todo, err
}

func GetTodoById(id int64) (*Todo, error) {
	var todo Todo
	ctx := context.Background()
	err := db.DB.NewSelect().Model(&todo).Where("id = ?", id).Scan(ctx)
	return &todo, err
}

func UpdateTodo(todo Todo) error {
	ctx := context.Background()
	var orig Todo
	orig.Done = todo.Done
	orig.Content = todo.Content
	_, err := db.DB.NewUpdate().Model(&orig).Where("id = ?", todo.ID).Exec(ctx)
	return err
}

func DeleteTodo(id int64) error {
	_, err := GetTodoById(id)
	ctx := context.Background()
	var orig Todo
	if err == nil {
		_, err = db.DB.NewDelete().Model(&orig).Where("id = ?", id).Exec(ctx)
	}
	return err
}