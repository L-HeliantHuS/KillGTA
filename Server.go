package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

var connects = make(map[string]net.Conn)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // 用于生成随机字符串

// RandStringRunes 随机生成字符串
func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// tcpServer 服务器主函数
func tcpServer(port int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	// 延迟关闭所有的客户端
	defer func() {
		for _, conn := range connects {
			_ = conn.Close()
		}
	}()

	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR] 创建服务器失败: %v", err))
	}

	log.Printf("[SUCCESS] 启动服务器成功, %s, 等待连接...\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("[WARNING] 有一个客户端连接失败了! 请注意检查! 当前连接数为: ", len(connects))
		}

		id := RandStringRunes(6)
		connects[id] = conn
		go worker(id, conn)
	}
}

// worker 分配给每个客户端的协程
func worker(id string, conn net.Conn) {
	log.Printf("[INFO] 新的客户端上线, ta的id为: %s, ta的地址为%s, 当前共有%d个客户端! \n", id, conn.RemoteAddr().String(), len(connects))

	for {
		buf := make([]byte, 4)
		readLen, err := conn.Read(buf)
		if err != nil {
			// 从列表中删除这个主机
			delete(connects, id)

			log.Println("[WARNING] ", conn.RemoteAddr().String(), "下线了, 当前在线主机数量: ", len(connects))

			break
		}

		msg := string(buf[:readLen])
		if msg == "kill" {
			// 收到任何一个人的kill指令 服务器都广播kill到各个客户端
			sendKills()
		} else if msg == "ping" {
			_, err := conn.Write([]byte("pong"))
			if err != nil {
				log.Println("[ERROR] ", conn.RemoteAddr().String(), "心跳测试响应失败.")
			}
		}
	}
}

// sendKills 发送给每个客户端Kill指令
func sendKills() {
	for _, conn := range connects {
		go func(conn net.Conn) {
			n, err := conn.Write([]byte("kill"))
			if err != nil {
				log.Println("[ERROR] 向", conn.RemoteAddr().String(), "发送Kill指令失败!")
			} else {
				log.Printf("[SUCCESS] 向%s发送Kill指令成功, 发送字节数%d \n", conn.RemoteAddr().String(), n)
			}
		}(conn)
	}
}

// status 启动一个协程来监测客户端数量
func status() {
	for {
		time.Sleep(time.Second * 30)

		log.Println("[SERVER-INFO] 当前在线主机: ", len(connects))
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 25155, "port")
	flag.Parse()
	go status()

	tcpServer(port)
}
