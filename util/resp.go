package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type H struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

// 返回失败
func RespFail(write http.ResponseWriter, msg string) {
	Resp(write, -1, nil, msg)
}

// 返回成功
func RespOk(write http.ResponseWriter, data interface{}, msg string) {
	Resp(write, 0, data, msg)
}

// 返回函数
func Resp(write http.ResponseWriter, code int, data interface{}, msg string) {
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)
	res := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ret, err := json.Marshal(res)
	if err != nil {
		log.Print(err.Error())
	}
	write.Write(ret)

}

func RespOKList(w http.ResponseWriter, lists interface{}, total interface{}) {
	// 分页
	RespList(w, 0, lists, total)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	// 将结构体转换为json字符串
	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(ret)
}
