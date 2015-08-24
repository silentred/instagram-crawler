package main

import (
	"encoding/json"
	"fmt"
	//"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	"time"
)

func getIndex(req *http.Request, params martini.Params, rd render.Render, session sessions.Session, db *sqlx.DB) {

	queries := req.URL.Query()
	fmt.Println(queries)

	// appName := params["appName"]
	// fmt.Println(appName)

	//session.Set("hello", "world")

	// You can also get a single result, a la QueryRow
	// pic := Picture{}
	// err := db.Get(&pic, "SELECT * FROM picture WHERE orig_id=1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%#v\n", pic)

	// user := &User{Id: "123", Name: "jason", AccessToken: "123123", LastAuthTime: 432123, Valid: 1}
	// ret := user.Insert()
	// fmt.Println(ret)
	user := &User{}
	user.GetById("123")
	fmt.Println(*user)

	url := fmt.Sprintf(API_URL, ClientId, RedirectUrl)

	rd.HTML(200, "auth_index", H{
		"url": url,
	})

}

func getRedirectBack(req *http.Request, rd render.Render, db *sqlx.DB, session sessions.Session) {
	var userId, userName, token string

	// request the token with code
	code := req.URL.Query().Get("code")
	fmt.Println(code)

	err_reason := req.URL.Query().Get("error_reason")
	if err_reason != "" {
		rd.Redirect("/oauth/instagram", 302)
	}

	// get username, userid
	postBody := map[string]string{
		"client_id":     ClientId,
		"client_secret": ClientSecret,
		"grant_type":    "authorization_code",
		"redirect_uri":  RedirectUrl,
		"code":          code,
	}
	content := HttpPost("https://api.instagram.com/oauth/access_token", postBody, nil, nil)
	fmt.Println(content)

	var mainContent map[string]*json.RawMessage
	json.Unmarshal([]byte(content), &mainContent)
	token = string(*mainContent["access_token"])

	var userInfo map[string]*json.RawMessage
	json.Unmarshal(*mainContent["user"], &userInfo)
	userId = string(*userInfo["id"])
	userName = string(*userInfo["username"])
	nowTime := int(time.Now().Unix())

	// save user, token.
	user := &User{Id: userId, Name: userName, AccessToken: token, LastAuthTime: nowTime, Valid: 1}
	_ = user.Insert()

	// @TODO Or update

	rd.HTML(200, "redirect_back", H{
		"userName": userName,
		"userId":   userId,
	})

}

var Jobs chan Picture
var Quit chan int
var NextURL chan string
var DoneCnt <-chan int

func init() {
	Jobs = make(chan Picture, 50)
	Quit = make(chan int)

	NextURL = make(chan string)
	//nextURL <- RecentURL // 通过传入 url来启动抓取

	// 任务队列放在全局执行，ws来控制是否开始或停止。 任务队列如果放在wsHandler中，一旦连接关闭，整个任务就停止了。
	out := preparePicture(Jobs, NextURL, Quit)
	DoneCnt = savingPicture(out)
}

func wsHandler(req *http.Request, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) {

	go func() {
		for {
			val := <-DoneCnt
			ret := &Message{Action: "DoneCnt", Data: string(val)}
			sender <- ret
		}
	}()

	for {
		select {
		case msg := <-receiver:
			// In your app, collect the senders of different clients and do something useful with them
			switch msg.Action {
			case "start":
				// get userId, get access_token from DB
				userId := req.URL.Query().Get("userId")
				user := &User{}
				user.GetById(userId)
				if user == nil {
					log.Fatal("cannot get user by id: " + userId)
				}
				access_token := user.AccessToken

				// give value to nextURL, so the goroutine could start
				url := fmt.Sprintf(RecentURL, access_token)
				NextURL <- url
			case "stop":
				Quit <- 1
			}

		case <-done:
			// the client disconnected, so you should return / break if the done channel gets sent a message
			fmt.Println("client close the connection. is server loop still going?")
			//return
		case err := <-errorChannel:
			// Uh oh, we received an error. This will happen before a close if the client did not disconnect regularly.
			// Maybe useful if you want to store statistics
			log.Fatal(err)
		}
	}

}

/**
 *  get pics and next_url from api, save next_url to channel
 */
func getPictureFromApi(url chan string) []Picture {
	pics := make([]Picture, 0, 30)
	next := <-url
	// get from api
	var result map[string]map[string]*json.RawMessage
	content := HttpGet(next, nil, nil)
	json.Unmarshal([]byte(content), &result)
	if errName := *result["meta"]["error_type"]; errName != nil {
		log.Fatal("error type: " + string(errName) + ", probably access_token expired")
	}

	imageId := *result["data"]["id"]
	imagesBytes := result["data"]["images"]
	var images map[string]map[string]*json.RawMessage
	json.Unmarshal(*imagesBytes, &images)
	imageUrl := *images["standard_resolution"]["url"]
	nowTime := int(time.Now().Unix())

	pic := Picture{Id: string(imageId), Url: string(imageUrl), Status: 0, CreatedTime: nowTime}
	pics = append(pics, pic)

	url <- string(*result["pagination"]["next_url"])
	return pics
}

func preparePicture(jobs chan Picture, url chan string, quit chan int) <-chan Picture {

	// continously getting pics from api, fill jobs with recieved pics
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				pics := getPictureFromApi(url)
				for _, pic := range pics {
					Jobs <- pic
				}
			}

		}
	}()
	return jobs
}

func savingPicture(jobs <-chan Picture) <-chan int {
	done := make(chan int, 30)
	// continously saving the pics from jobs , and downloading them
	go func() {
		for {
			pic := <-jobs
			pic.Insert()
			Download(pic.Url, DeterminDst(pic.Url))
			done <- 1
		}
	}()
	return done
}
