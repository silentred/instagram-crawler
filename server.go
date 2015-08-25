package main

import (
	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
)

const (
	session_key  = "9FsmWODVM90d5R9Wo4B6Fol8IXWhqH4m"
	API_URL      = "https://api.instagram.com/oauth/authorize/?client_id=%s&redirect_uri=%s&response_type=code"
	ClientId     = "875c559eb1964e308b840dfcac0533a2"
	ClientSecret = "b29d172d74ad42699e972df632a867fc"
	RedirectUrl  = "http://50.117.7.122:3000/oauth/instagram/back"

	RecentURL     = "https://api.instagram.com/v1/users/self/feed?access_token=%s"
	UserRecentURL = "https://api.instagram.com/v1/users/%s/media/recent?access_token=%s"
)

// db作为全局变量存在不太好，想办法从Injecter中取出比较好
var mysqlDB *sqlx.DB

type H map[string]interface{}

func main() {
	m := martini.Classic()
	// render html templates from templates directory
	configRender(m)

	// map DB to context
	mysqlDB = initMySQL(m)

	// serve static files; default is public
	m.Use(martini.Static("public"))

	// use session
	store := sessions.NewCookieStore([]byte(session_key))
	m.Use(sessions.Sessions("go_session", store))

	m.Group("/oauth", func(r martini.Router) {
		r.Get("/:appName", getIndex)
		r.Get("/instagram/back", getRedirectBack)
	})

	m.Get("/ws", sockets.JSON(Message{}), wsHandler)

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
	db, err := sqlx.Connect("mysql", "test:test@tcp(127.0.0.1:3306)/instagram")
	if err != nil {
		log.Fatalln(err)
	}

	m.Map(db)

	return db
}
