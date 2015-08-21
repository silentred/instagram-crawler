package main

import (
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
)

const (
	session_key = "9FsmWODVM90d5R9Wo4B6Fol8IXWhqH4m"
	ClientId    = "https://api.instagram.com/oauth/authorize/?client_id=%s&redirect_uri=%s&response_type=code"
)

type H map[string]interface{}

func main() {
	m := martini.Classic()
	// render html templates from templates directory
	configRender(m)

	// map DB to context
	_ = initMySQL(m)

	// serve static files; default is public
	m.Use(martini.Static("public"))

	// use session
	store := sessions.NewCookieStore([]byte(session_key))
	m.Use(sessions.Sessions("go_session", store))

	m.Group("/oauth", func(r martini.Router) {
		r.Get("/:appName", getIndex)
	})

	m.Run()
}

func configRender(m *martini.ClassicMartini) {
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                              // Specify what path to load the templates from.
		Layout:     "layout",                                 // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl.html", ".html", ".tmpl"}, // Specify extensions to load for templates.
		//Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		//Delims: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		//Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		//IndentJSON: true, // Output human readable JSON
		//IndentXML: true, // Output human readable XML
		//HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))
}

func initMySQL(m *martini.ClassicMartini) *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:@tcp(127.0.0.1:3306)/instagram")
	if err != nil {
		log.Fatalln(err)
	}

	m.Map(db)

	return db
}
