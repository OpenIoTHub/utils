package gateway

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	nettool "github.com/OpenIoTHub/utils/net"
	"github.com/OpenIoTHub/utils/net/p2p"
	"github.com/libp2p/go-yamux"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"time"
)

func NewP2PCtrlAsServer(stream net.Conn, ctrlmMsg *models.ReqNewP2PCtrlAsServer, token *models.TokenClaims) {
	//监听一个随机端口号，接受P2P方的连接
	externalUDPAddr, listener, err := p2p.GetP2PListener(token)
	if err != nil {
		log.Println(err)
		return
	}
	p2p.SendPackToPeerByReqNewP2PCtrlAsServer(listener, ctrlmMsg)
	//开始转kcp监听
	go kcpListener(listener, token)
	//TODO：发送认证码用于后续校验
	msg.WriteMsg(stream, &models.RemoteNetInfo{
		IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
		IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
		ExternalIp:   externalUDPAddr.IP.String(),
		ExternalPort: externalUDPAddr.Port,
	})
	//TODO:这里控制连接的处理？
	stream.Close()
}

//TODO：listener转kcp服务侦听
func kcpListener(listener *net.UDPConn, token *models.TokenClaims) {
	kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	kcplis.SetDeadline(time.Now().Add(time.Second * 5))
	//为了防范风险，只接受一个kcp请求
	//for {
	log.Println("start p2p kcp accpet")
	kcpconn, err := kcplis.AcceptKCP()
	if err != nil {
		kcplis.Close()
		if kcpconn != nil {
			kcpconn.Close()
		}
		log.Println(err.Error())
		return
	}
	//配置
	nettool.SetYamuxConn(kcpconn)

	log.Println("accpeted")
	log.Println(kcpconn.RemoteAddr())
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		kcplis.Close()
	}
	err = kcpConnHdl(kcpconn, token)
	if err != nil {
		kcplis.Close()
	}
	//}
}

func kcpConnHdl(kcpconn net.Conn, token *models.TokenClaims) error {
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*3)
	if err != nil {
		kcpconn.Close()
		log.Println(err.Error())
		return err
	}
	switch m := rawMsg.(type) {
	//TODO:初步使用ping、pong握手，下一步应该弄成验证校验身份
	case *models.Ping:
		{
			fmt.Printf("P2P握手ping")
			_ = m
			msg.WriteMsg(kcpconn, &models.Pong{})
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := yamux.Server(kcpconn, config)
			if err != nil {
				log.Println(err.Error())
			}
			go dlSubSession(session, token)
			fmt.Printf("Client作为Serverp2p打洞成功！")
			return nil
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
		kcpconn.Close()
		return fmt.Errorf("获取到了一个未知的P2P握手消息")
	}
}
