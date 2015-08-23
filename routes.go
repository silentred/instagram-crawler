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

	url := fmt.Sprintf(API_URL, ClientId, RedirectUrl)

	rd.HTML(200, "auth_index", H{
		"url": url,
	})

}

func getRedirectBack(req *http.Request, rd render.Render, db *sqlx.DB, session sessions.Session) {
	var userId, userName, token string

	// request the token with code
	code := req.URL.Query().Get("code")
	err_reason := req.URL.Query().Get("error_reason")
	if err_reason != "" {
		rd.Redirect("/oauth/instagram", 302)
	}

	// get username, userid
	params := map[string]string{
		"client_id":     ClientId,
		"client_secret": ClientSecret,
		"grant_type":    "authorization_code",
		"redirect_uri":  RedirectUrl,
		"code":          code,
	}
	content := HttpGet("https://api.instagram.com/oauth/access_token", params, nil)

	var mainContent map[string]*json.RawMessage
	json.Unmarshal([]byte(content), &mainContent)
	token = string(*mainContent["access_token"])

	var userInfo map[string]*json.RawMessage
	json.Unmarshal(*mainContent["user"], &userInfo)
	userId = string(*userInfo["id"])
	userName = string(*userInfo["username"])
	nowTime := int(time.Now().Unix())

	// save user, token.
	user := User{Id: userId, Name: userName, AccessToken: token, LastAuthTime: nowTime, Valid: 1}
	_ = user.Insert()

	// @TODO Or update

	rd.HTML(200, "redirect_back", H{
		"userName": userName,
	})

}

func init() {
	jobs := make(chan Picture, 50)
	quit := make(chan int)

	nextURL := make(chan string)
	nextURL <- RecentURL

	// 任务队列放在全局执行，ws来控制是否开始或停止。 任务队列如果放在wsHandler中，一旦连接关闭，整个任务就停止了。
	out := preparePicture(jobs, nextURL, quit)
	doneCnt := savingPicture(out)
	
}

func wsHandler(params martini.Params, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) {

	// ticker = time.After(30 * time.Minute)
	
	for {
		select {
		case msg := <-receiver:
			// here we simply echo the received message to the sender for demonstration purposes
			// In your app, collect the senders of different clients and do something useful with them
			switch msg.Action {
			case "start":
				

			case "stop":
				quit <- 1
			}

			go func () {
				for {
					val := <- doneCnt
					ret := &Message{Action: "doneCnt", Data: val}
					sender <- ret
				}
			}
			
		//case <-ticker:
		// This will close the connection after 30 minutes no matter what
		// To demonstrate use of the disconnect channel
		// You can use close codes according to RFC 6455
		//	disconnect <- websocket.CloseNormalClosure
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
	var result map[string]*json.RawMessage
	url <- string(*result["next_url"])
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
				// nextURL is a global channel
				pics := getPictureFromApi(url)
				for _, pic := range pics {
					jobs <- pic
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
			pic <- jobs
			pic.Insert()
			Download(pic.Url, DeterminDst(pic.Url))
			done <- 1
		}
	}()
	return done
}
