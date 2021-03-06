package src

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const YYYYMMDD = "2006-01-02"

type Reply struct {
	Code  int    `json:"code"`
	Datas Datas  `json:"datas"`
	Msg   string `json:"msg"`
}

type Datas struct {
	State   int     `json:"state"`
	Fields  []Field `json:"fields"`
	NowData string  `json:"now_data"`
}

type Field struct {
	FieldCode    string `json:"field_code"`
	UserSetValue string `json:"user_set_value"`
}

type Punch struct {
	PunchForm string `json:"punch_form"`
	Date      string `json:"date"`
}

type today struct {
	Date string `json:"date"`
}

func GetPunchForm(JSESSIONID *http.Cookie, nginx *http.Cookie) []Field {
	url := "https://yqtb.sut.edu.cn/getPunchForm"
	var currentDate today
	var reply Reply
	currentDate.Date = time.Now().UTC().Format(YYYYMMDD)
	//工大传统艺能今天打明天的卡,所以获取昨天的打卡记录其实是获取今天的，真的ybb
	DateJson, err := json.Marshal(currentDate)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(DateJson)
	req, _ := http.NewRequest("POST", url, reader)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(req.Body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
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
	return reply.Datas.Fields
}

// OnePunch 给我打卡,三回啊三回
// p.s.真不是我玩烂梗,而是这个龟儿子后端真的要提交三次才行,应该是一种验证方式
func OnePunch(fields []Field, JSESSIONID *http.Cookie, nginx *http.Cookie) bool {
	url := "https://yqtb.sut.edu.cn/punchForm"
	var (
		punch Punch
		reply Reply
	)

	//创建响应体map
	punchMap := make(map[string]string)
	for _, v := range fields {
		punchMap[v.FieldCode] = v.UserSetValue
	}
	jsonBytes, _ := json.Marshal(punchMap)
	punch.PunchForm = string(jsonBytes)
	punch.Date = time.Now().Add(time.Hour * 24).UTC().Format(YYYYMMDD)
	jsonBytes, _ = json.Marshal(punch)
	//注:学校那个打卡软件后端的作者肯定是个憨批,这个憨憨他非要把json给stringify了再解析(草
	//这样就导致了请求体就是个拼接了的字符串,简直就是一坨稀饭,而且最他喵离谱的是json变量还全是汉语拼音,就离谱
	//看得出来前端已经很努力在做解析了,搞了一堆方法给字符串加斜杠
	//p.s.我看源码的时候感觉前端都要蚌埠住了,我要是那前端我保证给做后台解析的那小子一个大耳刮子!
	body := bytes.NewReader(jsonBytes)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(req.Body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	req.AddCookie(JSESSIONID)
	req.AddCookie(nginx)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &reply)
	return reply.Code == 200
}
