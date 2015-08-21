package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
)

func getIndex(req *http.Request, params martini.Params, rd render.Render, session sessions.Session, db *sqlx.DB) {
	queries := req.URL.Query()
	fmt.Println(queries)

	appName := params["appName"]
	fmt.Println(appName)

	//session.Set("hello", "world")

	// You can also get a single result, a la QueryRow
	pic := Picture{}
	err := db.Get(&pic, "SELECT * FROM picture WHERE orig_id=1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", pic)

	rd.HTML(200, "auth_index", H{
		"title": "Hello Title",
		"name":  "Jason Bourn",
	})

}
