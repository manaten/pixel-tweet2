package pixeltweet

import (
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
)

var cached_templates = template.Must(template.ParseGlob("templates/*.html"))

type Entry struct {
	Title       string
	Image       string
	Description string
}

func init() {
	m := martini.Classic()
	m.Get("/", index)
	m.Group("/api/entry", func(r martini.Router) {
		// r.Get("/:id", getEntry)
		r.Post("", newEntry)
		// r.Put("/:id", updateEntry)
		// r.Delete("/:id", deleteEntry)
	})

	http.Handle("/", m)
}

func newEntry() string {
	return "{}"
}

func index(res http.ResponseWriter, req *http.Request) {
	e := []Entry{
		Entry{"hoge", "fuga", "hage"},
		Entry{"foo", "bar", "baz"},
	}
	cached_templates.ExecuteTemplate(res, "index.html", e)
}
