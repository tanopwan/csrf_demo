package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/", index)
	e.GET("/clickme", clickme)
	e.GET("/formpost", formpost)
	e.GET("/autopost", autopost)
	e.GET("/xhrjson", xhrjson)
	e.GET("/xhr", xhr)
	e.Logger.Fatal(e.Start(":1324"))
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func clickme(c echo.Context) error {
	return c.Render(http.StatusOK, "clickme", nil)
}

func formpost(c echo.Context) error {
	return c.Render(http.StatusOK, "formpost", nil)
}

func autopost(c echo.Context) error {
	return c.Render(http.StatusOK, "autopost", nil)
}

func xhrjson(c echo.Context) error {
	return c.Render(http.StatusOK, "xhrjson", nil)
}

func xhr(c echo.Context) error {
	return c.Render(http.StatusOK, "xhr", nil)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
