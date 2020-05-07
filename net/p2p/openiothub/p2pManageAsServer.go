package openiothub

import (
	"errors"
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

func MakeP2PSessionAsServer(stream net.Conn, TokenModel *models.TokenClaims) (*yamux.Session, error) {
	//TODO:这里控制连接的处理？
	if stream != nil {
		defer stream.Close()
	} else {
		return nil, errors.New("stream is nil")
	}
	//监听一个随机端口号，接受P2P方的连接
	ExternalUDPAddr, listener, err := p2p.GetP2PListener(TokenModel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = msg.WriteMsg(stream, &models.ReqNewP2PCtrlAsClient{
		IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
		IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
		ExternalIp:   ExternalUDPAddr.IP.String(),
		ExternalPort: ExternalUDPAddr.Port,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rawMsg, err := msg.ReadMsgWithTimeOut(stream, time.Second*5)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	switch m := rawMsg.(type) {
	case *net.UDPAddr:
		{
			p2p.SendPackToPeerByUDPAddr(listener, m)
			err = msg.WriteMsg(stream, &models.OK{})
			if err != nil {
				log.Println(err)
				return nil, err
			}
			log.Println("发送到p2p成功，等待连接")
			defer listener.Close()
			return kcpListener(listener.LocalAddr().(*net.UDPAddr))
		}
	default:
		log.Println("不是ReqNewP2PCtrlAsServer")
		return nil, errors.New("不是ReqNewP2PCtrlAsServer")
	}
}

//TODO：listener转kcp服务侦听
func kcpListener(laddr *net.UDPAddr) (*yamux.Session, error) {
	//kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
	kcplis, err := kcp.ListenWithOptions(laddr.String(), nil, 10, 3)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	kcplis.SetDeadline(time.Now().Add(time.Second * 5))
	//为了防范风险，只接受一个kcp请求
	//for {
	log.Println("start p2p kcp accpet")
	kcpconn, err := kcplis.AcceptKCP()
	if err != nil {
		if kcplis != nil {
			kcplis.Close()
		}
		if kcpconn != nil {
			kcpconn.Close()
		}
		log.Println(err.Error())
		return nil, err
	}
	//配置
	nettool.SetYamuxConn(kcpconn)
	log.Println("accpeted")
	log.Println(kcpconn.RemoteAddr())
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		log.Println(err)
		kcplis.Close()
	}
	return kcpConnHdl(kcpconn)
}

func kcpConnHdl(kcpconn *kcp.UDPSession) (*yamux.Session, error) {
	//:TODO 超时返回
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*5)
	if err != nil {
		kcpconn.Close()
		log.Println(err.Error())
		return nil, err
	}
	switch m := rawMsg.(type) {
	//TODO:初步使用ping、pong握手，下一步应该弄成验证校验身份
	case *models.Ping:
		{
			log.Println("P2P握手ping")
			_ = m
			config := yamux.DefaultConfig()
			//config.EnableKeepAlive = false
			session, err := yamux.Client(kcpconn, config)
			if err != nil {
				kcpconn.Close()
				log.Println("yamux.Client:", err.Error())
				return nil, err
			}
			return session, err
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
	}
	return nil, err
}
