package p2p

import (
	"github.com/OpenIoTHub/utils/models"
	nettool "github.com/OpenIoTHub/utils/net"
	"log"
	"net"
	"strconv"
	"time"
)

// GetDialIpPort 获取一个随机UDP Dial的内部ip，端口，外部ip端口
func GetDialIpPort(token *models.TokenClaims) (localAddr, externalAddr *net.UDPAddr, err error) {
	raddr, err := net.ResolveUDPAddr("udp", token.Host+":"+strconv.Itoa(token.UDPApiPort))
	//udpaddr, err := net.ResolveUDPAddr("udp", "tencent-shanghai-v1.host.nat-cloud.com:34321")
	if err != nil {
		return nil, nil, err
	}
	udpconn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}
	defer udpconn.Close()
	externalUDPAddr, err := nettool.GetExternalIpPortByUDP(udpconn, token)
	if err != nil {
		udpconn.Close()
		log.Println(err)
		return
	}
	//return strings.Split(udpconn.LocalAddr().String(), ":")[0]
	localAddr = udpconn.LocalAddr().(*net.UDPAddr)
	return localAddr, externalUDPAddr, err
}

// GetP2PListener 获取一个随机UDP Listen的内部ip，端口，外部ip端口
func GetP2PListener(token *models.TokenClaims) (externalUDPAddr *net.UDPAddr, listener *net.UDPConn, err error) {
	listener, err = net.ListenUDP("udp", nil)
	if err != nil {
		return
	}
	//获取监听的端口的外部ip和端口
	externalUDPAddr, err = nettool.GetExternalIpPortByUDP(listener, token)
	if err != nil {
		listener.Close()
		log.Println(err)
		return
	}
	return
}

// GetNewListener 把旧的Listener关闭创建一个新的Listener返回，本地地址相同
func GetNewListener(oldListener *net.UDPConn) (newListener *net.UDPConn, err error) {
	if oldListener != nil {
		oldListener.Close()
	}
	newListener, err = net.ListenUDP("udp", oldListener.LocalAddr().(*net.UDPAddr))
	return
}

// SendPackToPeerByUDPAddr client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeerByUDPAddr(listener *net.UDPConn, raddr *net.UDPAddr) {
	log.Println("发送包到远程：", raddr.IP, raddr.Port)
	//发送5次防止丢包，稳妥点
	for i := 1; i <= 5; i++ {
		listener.WriteToUDP([]byte("packFromPeer"), raddr)
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Millisecond * 200)
}

// SendPackToPeerByRemoteNetInfo client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeerByRemoteNetInfo(listener *net.UDPConn, ctrlmMsg *models.RemoteNetInfo) {
	log.Println("发送包到远程：", ctrlmMsg.ExternalIp, ctrlmMsg.ExternalPort)
	SendPackToPeerByUDPAddr(listener, &net.UDPAddr{
		IP:   net.ParseIP(ctrlmMsg.ExternalIp),
		Port: ctrlmMsg.ExternalPort,
	})
}

// SendPackToPeerByReqNewP2PCtrlAsServer client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeerByReqNewP2PCtrlAsServer(listener *net.UDPConn, ctrlmMsg *models.ReqNewP2PCtrlAsServer) {
	SendPackToPeerByRemoteNetInfo(listener, &models.RemoteNetInfo{
		IntranetIp:   ctrlmMsg.IntranetIp,
		IntranetPort: ctrlmMsg.IntranetPort,
		ExternalIp:   ctrlmMsg.ExternalIp,
		ExternalPort: ctrlmMsg.ExternalPort,
	})
}

// SendPackToPeerByReqNewP2PCtrlAsClient client通过指定listener发送数据到explorer指定的p2p请求地址
func SendPackToPeerByReqNewP2PCtrlAsClient(listener *net.UDPConn, ctrlmMsg *models.ReqNewP2PCtrlAsClient) {
	SendPackToPeerByRemoteNetInfo(listener, &models.RemoteNetInfo{
		IntranetIp:   ctrlmMsg.IntranetIp,
		IntranetPort: ctrlmMsg.IntranetPort,
		ExternalIp:   ctrlmMsg.ExternalIp,
		ExternalPort: ctrlmMsg.ExternalPort,
	})
}
