package ctrl

import (
	"fmt"
	"hao-im/util"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.MkdirAll("./mnt", os.ModePerm)
}
func Upload(w http.ResponseWriter, r *http.Request) {
	UploadLocal(w, r)
}

// 1.存储位置  ./mnt
// 2.url格式 /mnt/xxx.png  需要确保网络能访问/mnt
func UploadLocal(writer http.ResponseWriter, request *http.Request) {
	// todo 获得上传的文件源
	srcFile, head, err := request.FormFile("file")
	if err != nil {
		util.RespFail(writer, err.Error())
	}
	// 创建新文件
	suffix := ".png"
	// 如果前端文件包含后缀
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// 如果前端指定了filetype
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	// 生成文件ming
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dsfile, err := os.Create("./mnt/" + filename)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	// 将源文件copy到新文件
	_, err = io.Copy(dsfile, srcFile)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	// todo 将新文件路径转换为url地址
	url := "/mnt/" + filename
	// todo 响应到前端
	util.RespOk(writer, url, "")
}
