package main

import (
	"hao-im/ctrl"
	"html/template"
	"log"
	"net/http"
)

// 注册模板文件
func RegisterView() {
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil {
		// 打印并直接退出
		log.Fatal(err.Error())
	}
	for _, v := range tpl.Templates() {
		tplName := v.Name()
		http.HandleFunc(tplName, func(writer http.ResponseWriter, request *http.Request) {
			tpl.ExecuteTemplate(writer, tplName, nil)
		})

	}
}

// 注册首页自动跳转
func RegisterIndex() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		http.Redirect(writer, request, "/user/login.shtml", http.StatusFound)
	})
}

func main() {
	// 登录
	http.HandleFunc("/user/login", ctrl.UserLogin)
	// 注册
	http.HandleFunc("/user/register", ctrl.UserRegister)
	// 查找用户
	http.HandleFunc("/user/find",ctrl.FindUserById)


	// 加载个人全部好友
	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriend)
	// 加载个人全部群
	http.HandleFunc("/contact/loadcommunity",ctrl.LoadCommunity)
	// 加群
	http.HandleFunc("/contact/joincommunity",ctrl.JoinCommunity)
	// 创建群
	http.HandleFunc("/contact/createcommunity",ctrl.CreateCommunity)
	// 添加好友
	http.HandleFunc("/contact/addfriend", ctrl.AddFriend)

	// 聊天
	http.HandleFunc("/chat", ctrl.Chat)
	// 上传文件
	http.HandleFunc("/attach/upload", ctrl.Upload)



	// 注册静态文件
	http.Handle("/asset/", http.FileServer(http.Dir('.')))
	http.Handle("/mnt/", http.FileServer(http.Dir('.')))
	// 注册模板文件
	RegisterView()
	// 注册首页自动跳转
	RegisterIndex()
	// 启动web服务器
	http.ListenAndServe(":8000", nil)

}
