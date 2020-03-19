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
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes 随机生成字符串
func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func tcpServer(port int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	// 延迟关闭所有的客户端
	defer func() {
		for _, conn := range connects {
			conn.Close()
		}
	}()

	if err != nil {
		log.Fatal(fmt.Sprintf("创建服务器失败: %v", err))
	}

	fmt.Printf("[+] 启动服务器成功, %s, 等待连接...\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("有一个客户端连接失败了! 请注意检查! 当前连接数为: ", len(connects))
		}

		id := RandStringRunes(6)
		connects[id] = conn
		go worker(id, conn)
	}
}

func worker(id string, conn net.Conn) {
	fmt.Printf("[!]新的客户端上线, ta的id为: %s, ta的地址为%s, 当前共有%d个客户端! \n", id, conn.RemoteAddr().String(), len(connects))

	for {
		buf := make([]byte, 4)
		readLen, err := conn.Read(buf)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), "发送的消息读取失败, 他可能下线了.")

			// 从列表中删除这个主机
			delete(connects, id)

			break
		}

		msg := string(buf[:readLen])
		if msg == "kill" {
			// 收到任何一个人的kill指令 服务器都广播kill到各个客户端
			sendKills()
		} else if msg == "ping" {
			_, err := conn.Write([]byte("pong"))
			if err != nil {
				log.Println(conn.RemoteAddr().String(), "心跳测试响应失败.")
			}
		}
	}
}

func sendKills() {
	for _, conn := range connects {
		go func(conn net.Conn) {
			n, err := conn.Write([]byte("kill"))
			if err != nil {
				log.Println("向", conn.RemoteAddr().String(), "发送Kill指令失败!")
			} else {
				fmt.Printf("[+] 向%s发送Kill指令成功, 发送字节数%d \n", conn.RemoteAddr().String(), n)
			}
		}(conn)
	}
}

func status() {
	for {

		fmt.Println("[DEBUG] 当前在线主机: ", len(connects))

		time.Sleep(time.Second * 30)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 25155, "port")
	flag.Parse()
	go status()

	tcpServer(port)
}
