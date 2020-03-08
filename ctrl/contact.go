package ctrl

import (
	"hao-im/args"
	"hao-im/model"
	"hao-im/service"
	"hao-im/util"
	"net/http"
)

var contactService service.ContactService

// 加载个人好友
func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req, &arg)
	users := contactService.SearchFriend(arg.Userid)
	util.RespOKList(w, users, len(users))
}

// 加载全部群
func LoadCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req, &arg)
	communitys := contactService.SearchCommunity(arg.Userid)
	util.RespOKList(w, communitys, len(communitys))
}

// 添加群
func JoinCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req, &arg)
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	// 刷新用户的群组信息
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "")
	}
}

//添加好友
func AddFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req, &arg)
	// 调用service
	err := contactService.AddFriend(arg.Userid, arg.Dstid)

	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}

}

// 创建群
func CreateCommunity(w http.ResponseWriter, req *http.Request) {
	var arg model.Community
	util.Bind(req, &arg)
	com, err := contactService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, com, "")
	}

}
