package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	//user := &User{Id: "2", Name: "jason", AccessToken: "123123", LastAuthTime: 432123, Valid: 1}
	//ret := user.Insert()

	//expect(t, ret, 1)

	user := &User{}
	user.GetById("123")
	expect(t, user.Id, "123")
}

func TestPicture(t *testing.T) {

}
