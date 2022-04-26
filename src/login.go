package src

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// https://yqtb.sut.edu.cn/login
// method:POST
// {user_account:"181203501",user_password:"316919"}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Student struct {
	User_account  string `json:"user_account"`
	User_password string `json:"user_password"`
}

type login_reply struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Datas string `json:"datas"`
}

func loadStudentData(StudentBytes []byte) Student {
	var student Student
	_ = json.Unmarshal(StudentBytes, &student)
	return student
}

func initCookies(url string) (*http.Cookie, *http.Cookie) {
	var JSESSIONID, nginx *http.Cookie
	resp, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// for k, v := range resp.Header {
	// 	if k == "Set-Cookie" {
	// 		for key, value := range v {
	// 			if key != 0 {
	// 				nginx = value
	// 			} else {
	// 				JSESSIONID = value
	// 			}
	// 		}
	// 	}
	// }
	for k, v := range resp.Cookies() {
		if k != 0 {
			nginx = v
		} else {
			JSESSIONID = v
		}
	}
	return JSESSIONID, nginx
}

func DoLogin(jsonByte []byte) (bool, *http.Cookie, *http.Cookie) {
	url := "https://yqtb.sut.edu.cn"
	JSESSIONID, nginx := initCookies(url)

	url = fmt.Sprintf("%s/login", url)
	var reply login_reply

	reader := bytes.NewReader(jsonByte)
	req, err := http.NewRequest("POST", url, reader) //搞个新请求来
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(JSESSIONID)
	req.AddCookie(nginx)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(respBytes, &reply)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp.Cookies() {
		JSESSIONID = v
	}
	return reply.Code == 200, JSESSIONID, nginx
}
