package main

import (
	"encoding/json"
	. "gd-one-punch/src"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PunchNeededList() [][]byte {
	var students []Student = make([]Student, 6)
	var list [][]byte = make([][]byte, 6)
	//在这里填入一宿舍6个懒蛋的学号密码
	students[0] = Student{User_account: "111111111", User_password: "123456"}
	for i := 0; i < 6; i++ {
		if students[i].User_account == "" {
			break
		}
		list[i], _ = json.Marshal(students[i])
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
				return
			}
			success, JSESSIONID, nginx := DoLogin(v)
			err := json.Unmarshal(v, &student)
			if err != nil {
				log.Println("[+] 打卡完成!")
				return
			}
			if !success {
				log.Printf("[-] WARNING! 倒霉蛋%s登陆失败!", student.User_account)
				return
			}
			if !OnePunch(GetPunchForm(JSESSIONID, nginx), JSESSIONID, nginx) {
				log.Printf("[-] WARNING! 倒霉蛋%s打卡失败!", student.User_account)
				return
			}
			log.Printf("[+] %s 打卡成功!", student.User_account)
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
