package routes

import (
	"sample/models"

	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Data struct {
	Todos  []models.Todo
	Errors []error
}

func customFunc(todo *models.Todo) func([]string) []error {
	return func(values []string) []error {
		if len(values) == 0 || values[0] == "" {
			return nil
		}
		dt, err := time.Parse("2006-01-02T15:04 MST", values[0]+" JST")
		if err != nil {
			return []error{echo.NewBindingError("until", values[0:1], "failed to decode time", err)}
		}
		todo.Until = dt
		return nil
	}
}

func get(c echo.Context) error {
		todos, err := models.GetAllTodos()
		if err != nil {
			c.Logger().Error(err)
			return c.Render(http.StatusBadRequest, "index", Data{
				Errors: []error{errors.New("Cannot get todos")},
			})
		}
		return c.Render(http.StatusOK, "index", Data{Todos: todos})
}

func getRequest(c echo.Context) error {
		todos, err := models.GetAllTodos()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "TODOを取得できませんでした。")
		}
		return c.JSON(http.StatusOK, todos)
}

func post(c echo.Context) error {
		var todo models.Todo
		// フォームパラメータをフィールドにバインド
		errs := echo.FormFieldBinder(c).
			Int64("id", &todo.ID).
			String("content", &todo.Content).
			Bool("done", &todo.Done).
			CustomFunc("until", customFunc(&todo)).
			BindErrors()
		if errs != nil {
			return c.Render(http.StatusBadRequest, "index", Data{Errors: errs})
		} else if todo.ID == 0 {
			// ID が 0 の時は登録
			if todo.Content == "" {
				err := errors.New("Todo not found")
				c.Logger().Error(err)
			} else {
				_, err := models.CreateTodo(todo)
				if err != nil {
					c.Logger().Error(err)
					err = errors.New("Cannot update")
				}
			}
		} else {
			if c.FormValue("delete") != "" {
				// 削除
				err := models.DeleteTodo(todo.ID)
				if err != nil {
					c.Logger().Error(err)
					err = errors.New("Cannot update")
					return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
				}
			} else {
				// 更新
				err := models.UpdateTodo(todo)
				if err != nil {
					c.Logger().Error(err)
					err = errors.New("Cannot update")
					return c.Render(http.StatusBadRequest, "index", Data{Errors: []error{err}})
				}
			}

		}
		return c.Redirect(http.StatusFound, "/")
	}
