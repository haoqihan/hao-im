package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

const (
	CMD_SINGLE_MSG = 10 //点对点单聊，dstid是用户id
	CMD_ROOM_MSG   = 11 // 群聊信息，dstid是群id
	CMD_HEART      = 0  // 心跳信息，不做处理
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     //对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         //预览图片
	Url     string `json:"url,omitempty" form:"url"`         //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   //其他和数字相关的
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系表
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	//todo 检验接入是否合法
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isvalida := checkToken(userId, token)
	fmt.Println(id, token, isvalida)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//todo 获得conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// todo 获取用户全部群id
	comIds := contactService.SearchCommunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}
	// todo userid 和node形成绑定关系
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	// todo 完成发送逻辑 con
	go sendproc(node)
	// todo 完成接收逻辑、
	go recvproc(node)
	// 发送信息
	sendMsg(userId, []byte("hello world"))

}

// ws 接收协程
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// todo 对data进一步处理
		dispatch(data)
		// 把消息广播到局域网
		broadMsg(data)
		fmt.Printf("%s", data)
	}
}

//用来存放发送的要广播的数据
var udpsendchan = make(chan []byte, 1024)

//todo 将消息广播到局域网
func broadMsg(data []byte) {
	udpsendchan <- data

}
func init() {
	go udpsendproc()
	go udprecvproc()
}

//todo 完成udp数据的发送协程
func udpsendproc() {
	log.Println("start udpsendproc")
	// todo 使用udp协议拨号
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// todo 通过的到con发送信息
	for {
		select {
		case data := <-udpsendchan:
			_, err := con.Write(data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// todo 完成udp接收处理功能
func udprecvproc() {
	log.Println("start udprecvproc")
	// todo 监听udp广播端口
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
	}
	// todo 处理端口发过来的数据
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	// todo 解析data为message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// todo 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		// todo 群聊转发逻辑
		for userid, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				//自己排除，不发送
				if msg.Userid != userid {
					v.DataQueue <- data
				}
			}

		}
	case CMD_HEART:
		// todo 心跳 一遍多不做
	}

}

// 发送协程
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// todo 发送信息
func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}

}

// 检测token是否有效
func checkToken(userId int64, token string) bool {
	// 从数据库里查询并比对
	user := userService.Find(userId)
	return user.Token == token

}

//todo 添加新的群ID到用户的groupset中
func AddGroupId(userId int64, gid int64) {
	// 取得锁
	rwLocker.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	rwLocker.Unlock()
}
