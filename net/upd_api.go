package nettool

import (
	"fmt"
	"github.com/OpenIoTHub/utils/crypto"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func RunApiServer(port int) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		log.Println(err)
		return
	}
	go udpListener(listener)
}

func udpListener(listener *net.UDPConn) {
	data := make([]byte, 256)
	for {
		_, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		//fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		//:TODO 防阻塞
		go func() {
			_, err = listener.WriteToUDP([]byte(remoteAddr.String()), remoteAddr)
			if err != nil {
				fmt.Printf(err.Error())
			}
		}()
	}
}

//获取一个listener的外部地址和端口
func GetExternalIpPort(listener *net.UDPConn, token *crypto.TokenClaims) (ip string, port int, err error) {
	udpaddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.P2PApiPort))
	//udpaddr, err := net.ResolveUDPAddr("udp", "tencent-shanghai-v1.host.nat-cloud.com:34321")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	err = listener.SetDeadline(time.Now().Add(time.Duration(3 * time.Second)))
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	listener.WriteToUDP([]byte("getIpPort"), udpaddr)

	log.Println("发送到服务器确定成功！等待确定外网ip和port")
	data := make([]byte, 256)
	n, _, err := listener.ReadFromUDP(data)
	log.Println("获取api的UDP包成功，开始解析自己listener出口地址和端口")
	if err != nil {
		fmt.Printf("获取listener的出口出错: %s", err.Error())
		return "", 0, err
	}
	ipPort := string(data[:n])
	ip = strings.Split(ipPort, ":")[0]
	port, err = strconv.Atoi(strings.Split(ipPort, ":")[1])
	if err != nil {
		fmt.Printf(err.Error())
		log.Println("解析listener外部出口信息错误")
		return "", 0, err
	}

	err = listener.SetDeadline(time.Now().Add(time.Duration(99999 * time.Hour)))
	if err != nil {
		fmt.Printf("%s", err.Error())
		return "", 0, err
	}

	log.Println("我的公网IP:", strings.Split(ipPort, ":")[0])
	log.Println("内网的的出口端口:", port)
	return ip, port, err
}