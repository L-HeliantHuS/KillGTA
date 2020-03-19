package main

import (
	"flag"
	"fmt"
	"github.com/go-vgo/robotgo"
	"log"
	"net"
	"os/exec"
	"time"
)

var runStatus bool
var ipaddr string
var tcpServerConn net.Conn

func KillProcess() {
	killGTA := exec.Command("taskkill", "/F", "/T", "/IM", "GTA5.exe").Run()
	killRockstar := exec.Command("taskkill", "/F", "/T", "/IM", "launcher.exe").Run()
	if killGTA != nil {
		log.Println("未成功杀死GTA, 请检查GTA是否运行中. Error: ", killGTA.Error())
	} else {
		log.Println("成功杀死GTA! 检查你的首脑进度是否还存在吧！")
	}
	if killRockstar != nil {
		log.Println("未成功杀死R* Client, 请检查R* Client是否运行中. Error: ", killRockstar.Error())
	} else {
		log.Println("成功杀死R* Client! 检查你的首脑进度是否还存在吧！")
	}
}

// dialTcp 与服务器建立连接
func dialTcp(ip string, retryFlag bool) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		if !retryFlag {
			log.Fatal("连接服务器失败:", err)
		} else {
			log.Println("服务器重新失败. 还会进行尝试.")
			return
		}
	}


	fmt.Println("连接服务器成功, 服务器地址为:", conn.RemoteAddr().String())
	tcpServerConn = conn
}

// clientWorker 监听Server发送过来的Kill
func clientWorker() {
	runStatus = true
	for {
		buf := make([]byte, 4)
		n, err := tcpServerConn.Read(buf)
		if err != nil {
			runStatus = false
			break
		}
		msg := string(buf[:n])
		if msg == "kill" {
			go KillProcess()
		} else if msg == "ping" {
			log.Println("服务器在线, 而且在正常工作!")
		}
	}
}

// serverHeartTest 与服务器进行心跳测试
func serverHeartTest() {
	for {
		_, err := tcpServerConn.Write([]byte("ping"))
		if err != nil {
			log.Println("心跳包发送失败！请检查与服务器的连接.")
			dialTcp(ipaddr, true)
		} else {
			log.Println("[DEBUG] Server is onlined!")
		}
		if runStatus == false {
			go clientWorker()
		}
		time.Sleep(time.Second * 30)
	}
}

func main() {
	var online bool
	fmt.Println("程序运行时, 不要点击黑框框内部, 否则程序会暂停。")
	fmt.Println("杀死GTA的快捷键为F4, 请尝试点击看看是否有作用, 然后直接最小化本程序就可以了.")
	flag.BoolVar(&online, "online", true, "online")
	flag.Parse()

	if online {
		fmt.Println("输入服务器地址(列如 127.0.0.1:25155), 一定一定按照列子的格式填写！")
		fmt.Print(">>> ")
		_, err := fmt.Scanln(&ipaddr)

		if err != nil {
			log.Fatal("读取输入失败!")
		}

		fmt.Printf("读取输入成功: %s \n", ipaddr)

		dialTcp(ipaddr, false)
		defer tcpServerConn.Close()

		go clientWorker()
		go serverHeartTest()

		for {
			clickedF4 := robotgo.AddEvents("f4")
			if clickedF4 {
				go func() {
					write, err := tcpServerConn.Write([]byte("kill"))
					if err != nil {
						log.Println("[-] 向服务器发送kill失败")
					} else {
						log.Println("[+] 向服务器发送kill成功, 发送字节数为", write)
					}
				}()

				// 如果与服务器的连接突然断开了. 可以断掉自己的进程.
				if runStatus == false {
					go KillProcess()
				}
			}

			// 防止Hook卡死
			time.Sleep(2 * time.Second)
		}

	} else {
		fmt.Println("[!]以单机模式运行中...")
		// 单机模式运行 (一般没啥用, 不过可以秒关GTA= =)
		for {
			clickedF4 := robotgo.AddEvents("f4")
			if clickedF4 {
				KillProcess()

			}
		}
	}

}
