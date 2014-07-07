package pixeltweet

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

var cached_templates = template.Must(template.ParseGlob("templates/*.html"))

type Entry struct {
	Id          int64 `datastore:"-"`
	Title       string
	Image       string
	Description string
	Date        time.Time
}

func init() {
	m := martini.Classic()
	m.Get("/", index)
	m.Group("/api/entry", func(r martini.Router) {
		r.Get("", getEntry)
		r.Post("", newEntry)
		// r.Put("/:id", updateEntry)
		// r.Delete("/:id", deleteEntry)
	})

	http.Handle("/", m)
}

func entryKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Entry", "default_entry", 0, nil)
}

// datastoreからクエリで取得しJSONレスポンスで返すとこまでまとめていいかも
func getEntry(res http.ResponseWriter, req *http.Request) string {
	c := appengine.NewContext(req)
	q := datastore.NewQuery("Entry").Order("Date").Limit(50)
	entries := make([]Entry, 0, 50)
	keys, err := q.GetAll(c, &entries)
	if err != nil {
		return "ng"
	}

	for i, key := range keys {
		entries[i].Id = key.IntID()
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(entries)
	if err != nil {
		return "ng"
	}
	return string(json)
}

// 同じく、datastoreに突っ込んでレスポンス返すとこまでまとめていいかも
func newEntry(res http.ResponseWriter, req *http.Request) string {
	c := appengine.NewContext(req)
	e := Entry{
		0,
		req.FormValue("title"),
		req.FormValue("image"),
		req.FormValue("description"),
		time.Now(),
	}

	key := datastore.NewIncompleteKey(c, "Entry", entryKey(c))
	_, err := datastore.Put(c, key, &e)
	if err != nil {
		return "ng"
	}
	return "ok"
}

func index(res http.ResponseWriter, req *http.Request) {
	e := []Entry{
		Entry{0, "hoge", "fuga", "hage", time.Now()},
		Entry{0, "foo", "bar", "baz", time.Now()},
	}
	cached_templates.ExecuteTemplate(res, "index.html", e)
}
