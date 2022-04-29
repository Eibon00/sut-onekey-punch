package main

import (
	"encoding/json"
	"errors"
	"fmt"
	. "gd-one-punch/src"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const punchFileName = "punch.json"

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../")
}

func FileExists(filepath string) bool {
	stat, err := os.Stat(filepath)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !stat.IsDir()
}

func PunchNeededList() [][]byte {

	var list [][]byte
	punchFile := fmt.Sprintf("%s/%s", GetRootPath(), punchFileName)
	if FileExists(punchFile) {
		var punchConfig PunchFile

		jsonBytes, err := ioutil.ReadFile(punchFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonBytes, &punchConfig)
		for _, v := range punchConfig.Students {
			studentByte, _ := json.Marshal(v)
			list = append(list, studentByte)
		}
	} else {
		//在这里填入一宿舍6个懒蛋的学号密码
		var students []Student = make([]Student, 6)
		students[0] = Student{UserAccount: "111111111", UserPassword: "114514"}
		students[1] = Student{UserAccount: "222222222", UserPassword: "191981"}
		for i := 0; i < 6; i++ {
			if students[i].UserAccount == "" {
				break
			}
			studentByte, _ := json.Marshal(students[i])
			list = append(list, studentByte)
		}
	}
	return list
}

func main() {
	log.Println("[+] 开始打卡")
	ch := make(chan []byte, 6)

	go func() {
		var student Student
		for {
			v, ok := <-ch
			if !ok {
				log.Println("[+] 打卡完成!")
				return
			}
			success, JSESSIONID, nginx := DoLogin(v)
			err := json.Unmarshal(v, &student)
			if err != nil {

				return
			}

			if !success {
				log.Printf("[-] WARNING! 倒霉蛋%s登陆失败!", student.UserAccount)
			}

			if !OnePunch(GetPunchForm(JSESSIONID, nginx), JSESSIONID, nginx) {
				log.Printf("[-] WARNING! 倒霉蛋%s打卡失败!", student.UserAccount)
			} else {
				log.Printf("[+] %s 打卡成功!", student.UserAccount)
			}
		}

	}()

	for k, v := range PunchNeededList() {
		if v == nil {
			log.Printf("[-] data in list[%d] is empty,ignore...", k)
		} else {
			//每隔5~10分钟塞一个参数进来
			latency := rand.Intn(5) + 5
			time.Sleep(time.Duration(latency) * time.Minute)
			ch <- v
		}
	}
	close(ch)
	time.Sleep(10 * time.Second)
}
