package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DeterminDst(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	return os.TempDir() + u.RequestURI()
}

func Download(url, dst string) int64 {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		//index := strings.LastIndex(dst, "/")
		//path := dst[:index]
		path, _ := filepath.Split(dst)
		os.MkdirAll(path, 0777)
		out, err := os.Create(dst)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		n, err := io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return n
	}
	return 0
}

// send http request GET
func HttpGet(url string, params, headers map[string]string) string {
	queryArray := make([]string, 0)
	for key, value := range params {
		queryArray = append(queryArray, key+"="+value)
	}
	queryString := strings.Join(queryArray, "&")
	fullUrl := url + "?" + queryString
	//fmt.Println(fullUrl)

	client := &http.Client{}
	request, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)

}

func HttpPost(apiUrl string, postBody, params, headers map[string]string) string {
	queryArray := make([]string, 0)
	for key, value := range params {
		queryArray = append(queryArray, key+"="+value)
	}
	queryString := strings.Join(queryArray, "&")
	fullUrl := apiUrl + "?" + queryString
	//fmt.Println(fullUrl)

	body := url.Values{}
	for key, value := range postBody {
		body.Add(key, value)
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", fullUrl, bytes.NewBufferString(body.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)

}
