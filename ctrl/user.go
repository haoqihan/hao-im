package ctrl

import (
	"fmt"
	"hao-im/model"
	"hao-im/service"
	"hao-im/util"
	"math/rand"
	"net/http"
)

var userService service.UserService

// 登录
func UserLogin(write http.ResponseWriter, request *http.Request) {
	// 解析参数
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(write, err.Error())
	} else {
		util.RespOk(write, user, "")
	}
}

// 注册
func UserRegister(writer http.ResponseWriter, request *http.Request) {
	// 解析参数
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	plainpwd := request.PostForm.Get("passwd")
	nickName := fmt.Sprintf("user%06d", rand.Int31())
	acatar := ""
	sex := model.SEX_UNKNOW
	user, err := userService.Register(mobile, plainpwd, nickName, acatar, sex)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}

}

// 查找某个用户
func FindUserById(w http.ResponseWriter, req *http.Request) {
	var user model.User
	util.Bind(req, &user)
	user = userService.Find(user.Id)
	if user.Id == 0 {
		util.RespFail(w, "该用户不存在")
	} else {
		util.RespOk(w, user, "")
	}
}
