package main

import (
	"log"
)

type Picture struct {
	Id          string `db:"orig_id"`
	Url         string `db:"pic_url"`
	Status      int    `db:"status"`
	CreatedTime int    `db:"created_time"`
}

type User struct {
	Id           string `db:"orig_id"`
	Name         string `db:"name"`
	AccessToken  string `db:"access_token"`
	LastAuthTime int    `db:"last_auth_time"`
	Valid        int    `db:"valid"`
}

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func (user User) Insert() int {
	sql := "insert into user (orig_id, name, access_token, last_auth_time, valid) values (:orig_id, :name, :access_token, :last_auth_time, :valid)"
	result, err := mysqlDB.NamedExec(sql, &user)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return int(rowsAffected)
}
