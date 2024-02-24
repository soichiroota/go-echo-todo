package main

import (
	"sample/db"
	"sample/models"
	"sample/routes"

	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func formatDateTime(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02 15:04")
}

func main() {
  db.Init()

	err := models.NewCreateTodoTable()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	routes.RegisterRoutes(e)

	e.Renderer = &Template{
		templates: template.Must(template.New("").
			Funcs(template.FuncMap{
				"FormatDateTime": formatDateTime,
			}).ParseFS(templates, "templates/*")),
	}

	staticFs, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FileSystem(http.FS(staticFs)))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fileServer)))
	e.Logger.Fatal(e.Start(":8989"))
}