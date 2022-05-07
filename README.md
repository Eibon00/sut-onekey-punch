# sut-onekey-punch


*沈阳工业大学自动每日健康打卡，~~让你和你的懒蛋室友解放双手~~*

***注意:本程序仅供学习交流使用,不要谎报瞒报,若影响防疫政策实施,本人概不负责。***


## 使用方法1：

登录自己的服务器（或PC\软路由\电视机顶盒等支持golang的设备）, 输入
> $ git clone https://github.com/Eibon00/sut-onekey-punch.git

将项目下载到本地

打开`main.go`文件, 找到如下代码:

````
students[0] = Student{User_account: "111111111", User_password: "123456"}
````

分别在`User_account`和`User_password`中填入自己的`学号`和`打卡密码`,\
可以添加多个`students[]`,最多再添加5个,从`student[1]`到`student[5]`

### 编译程序

> $ ./go build main.go

## 使用方法2：

**1.下载对应架构编译好的可执行文件**

**2.修改punch.json,并放在与可执行文件相同目录下**

比如可执行文件路径为`/home/eibon00/sut-onekey-punch/onekey-punch` \
则`punch.json`文件的路径应为`/home/eibon00/sut-onekey-punch/punch.json`

**文件大致内容如下**

```json
{
  "students": [
    {
      "user_account": "111111111",
      "user_password": "111111"
    },
    {
      "user_account": "222222222",
      "user_password": "222222"
    }
  ]
}
```

*p.s.程序会优先识别相同目录下的`punch.json`文件, 反而源码中的结构体优先级较低, 若该文件不存在,则加载程序中的结构体,但二者不能同时使用* \
*虽然由于使用了json文件定义的用户数量可以超过6个, 但仍不建议使用6个以上的json对象, 可能会产生不可预测的问题(懒得修bug)*

## 添加为定时运行

> $ crontab -e

#### 按`i`输入

```
0 3 * * * /home/你的用户名/sut-onekey-punch/main
```

*p.s.这样就可以每天上午11点运行一次了*

#### 或者

```
0 3 * * * /home/你的用户名/sut-onekey-punch/main >> /home/你的用户名/punch.log 2>&1
```

#### 按`ESC`输入`:wq`保存退出

*p.s.这样就可以将每天运行的结果输出到`punch.log`文件,方便查看*

***然后就能愉快使用啦~***