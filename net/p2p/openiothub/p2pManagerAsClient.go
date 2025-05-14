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

// 作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, TokenModel *models.TokenClaims) (p2pSubSession *yamux.Session, listener *net.UDPConn, err error) {
	ExternalUDPAddr, listener, err := p2p.GetP2PListener(TokenModel)
	if err != nil {
		log.Println(err.Error())
		stream.Close()
		if listener != nil {
			listener.Close()
		}
		return
	}
	msgsd := &models.ReqNewP2PCtrlAsServer{
		IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
		IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
		ExternalIp:   ExternalUDPAddr.IP.String(),
		ExternalPort: ExternalUDPAddr.Port,
	}
	err = msg.WriteMsg(stream, msgsd)
	if err != nil {
		log.Println(err)
		stream.Close()
		listener.Close()
		return
	}
	var rawMsg models.Message
	rawMsg, err = msg.ReadMsg(stream)
	if err != nil {
		log.Println(err)
		stream.Close()
		listener.Close()
		return
	}
	switch m := rawMsg.(type) {
	case *net.UDPAddr:
		{
			log.Println("remote net info")
			log.Println("===", m.String())
			//TODO:认证；同内网直连；抽象出公共函数？
			var kcpconn *kcp.UDPSession
			kcpconn, err = kcp.NewConn(m.String(), nil, 10, 3, listener)
			//设置
			if err != nil {
				log.Println(err.Error())
				stream.Close()
				listener.Close()
				if kcpconn != nil {
					kcpconn.Close()
				}
				return
			}
			nettool.SetYamuxConn(kcpconn)
			time.Sleep(time.Second)
			err = msg.WriteMsg(kcpconn, &models.Ping{})
			if err != nil {
				stream.Close()
				kcpconn.Close()
				listener.Close()
				log.Println(err)
				return
			}
			//:TODO 超时返回
			//var rawMsg models.Message
			rawMsg, err = msg.ReadMsgWithTimeOut(kcpconn, time.Second*5)
			if err != nil {
				stream.Close()
				kcpconn.Close()
				listener.Close()
				log.Println(err)
				return
			}
			switch m := rawMsg.(type) {
			case *models.Pong:
				{
					log.Println("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					config := yamux.DefaultConfig()
					//config.EnableKeepAlive = false
					p2pSubSession, err = yamux.Client(kcpconn, config)
					if err != nil {
						stream.Close()
						kcpconn.Close()
						listener.Close()
						log.Println("create sub session err:" + err.Error())
						return nil, nil, err
					}
					return
				}
			default:
				stream.Close()
				kcpconn.Close()
				listener.Close()
				log.Println("type err")
			}
		}
	default:
		stream.Close()
		listener.Close()
		log.Println("type err")
		err = errors.New("message不匹配")
		return
	}
	err = errors.New("没有匹配到对方发送过来的UDP地址")
	return
}
