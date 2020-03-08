package service

import (
	"errors"
	"fmt"
	"hao-im/model"
	"hao-im/util"
	"math/rand"
	"time"
)

type UserService struct {
}

// 注册函数
func (u *UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	/*
		mobile 手机号
		plainpwd 密码
		nickname 昵称
		avatar 头像
	*/
	// 检测手机号是否存在
	tmp := model.User{}
	_, err = DbEngin.Where("mobile = ?", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	// 如果存在返回已注册
	if tmp.Id > 0 {
		return tmp, errors.New("该手机号已经注册")
	}
	// 拼接参数
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	// 密码 （MD5 加密）
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()
	// token
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())

	// 插入信息
	_, err = DbEngin.InsertOne(&tmp)

	return tmp, nil
}

// 登录函数
func (u *UserService) Login(mobile, plainpwd string) (user model.User, err error) {
	// 通过手机号查询用户
	tmp := model.User{}
	_, err = DbEngin.Where("mobile = ?", mobile).Get(&tmp)
	if tmp.Id == 0 {
		return tmp, errors.New("该用户不存在")
	}
	// 查询到了对比密码
	if ! util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd) {
		return tmp, errors.New("密码不正确")
	}
	// 刷新token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	// 返回数据
	_, err = DbEngin.Where("id = ?", tmp.Id).Cols("token").Update(&tmp)
	return tmp, err
}

// 查找某个用户
func (s *UserService) Find(userId int64) (user model.User) {
	tmp := model.User{}
	DbEngin.ID(userId).Get(&tmp)
	return tmp

}
