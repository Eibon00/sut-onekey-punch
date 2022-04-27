package src

import (
	"bytes"
	"crypto/tls"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Student struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

type loginReply struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Datas string `json:"datas"`
}

func initCookies(url string) (*http.Cookie, *http.Cookie) {
	var JSESSIONID, nginx *http.Cookie
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //解决部分Linux操作系统上由于无法验证证书导致的panic
	}
	client := &http.Client{Transport: tr}
	resp, _ := client.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

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
	JSESSIONID, nginx := initCookies(url) //初始化cookies

	url = fmt.Sprintf("%s/login", url)
	var reply loginReply

	//直接修改默认http客户端,好孩子不要学
	//http.DefaultClient.Transport = &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		InsecureSkipVerify: true,
	//	},
	//}
	reader := bytes.NewReader(jsonByte)
	req, _ := http.NewRequest("POST", url, reader) //只要我不做异常处理,就是没有异常(确信
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(req.Body)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(JSESSIONID)
	req.AddCookie(nginx)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	resp, _ := c.Do(req)
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
