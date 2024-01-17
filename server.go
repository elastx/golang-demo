package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("pages/%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Static("/images", "images")
	e.Static("/css", "css")
	e.Static("/pages", "pages")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"hostname": getHostname(),
		})
	}).Name = "mainPage"

	e.Logger.Fatal(e.Start(":80"))
}
