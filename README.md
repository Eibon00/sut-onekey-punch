# sut-onekey-punch
### 沈阳工业大学自动每日健康打卡，让你~~和你的懒蛋室友~~解放双手
## 使用方法：
登录自己的服务器（或PC\软路由\电视机顶盒等支持golang的设备）, 输入
>$ git clone https://github.com/Eibon00/sut-onekey-punch.git 

将项目下载到本地

打开`main.go`文件, 找到如下代码:
````
var list [][]byte = make([][]byte, 6)
students[0] = Student{User_account: "111111111", User_password: "123456"}
````

分别在`User_account`和`User_password`中填入自己的`学号`和`打卡密码`,\
可以添加多个`students[]`,最多再添加5个,从`student[1]`到`student[5]`

### 编译程序
>$ ./go build main.go

### 添加为定时运行
>$ crontab -e

#### 按`i`输入 `0 11 * * * /home/你的用户名/sut-onekey-punch/main`
#### 按`ESC`输入`:wq`保存退出
*p.s.这样就可以每天上午11点运行一次了*
## 然后就能愉快使用啦~