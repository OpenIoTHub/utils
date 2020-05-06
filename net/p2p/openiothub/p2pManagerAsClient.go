package openiothub

import (
	"errors"
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

//作为客户端主动去连接内网client的方式创建穿透连接
func MakeP2PSessionAsClient(stream net.Conn, TokenModel *models.TokenClaims) (*yamux.Session, error) {
	if stream != nil {
		defer stream.Close()
	} else {
		return nil, errors.New("stream is nil")
	}
	localAddr, externalUDPAddr, err := p2p.GetDialIpPort(TokenModel)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	msgsd := &models.ReqNewP2PCtrlAsServer{
		IntranetIp:   localAddr.IP.String(),
		IntranetPort: localAddr.Port,
		ExternalIp:   externalUDPAddr.IP.String(),
		ExternalPort: externalUDPAddr.Port,
	}
	err = msg.WriteMsg(stream, msgsd)
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
			log.Println("remote net info")
			udpconn, err := net.DialUDP("udp", localAddr, m)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			//TODO:认证；同内网直连；抽象出公共函数？
			kcpconn, err := kcp.NewConn(fmt.Sprintf("%s:%d", m.IP, m.Port), nil, 10, 3, udpconn)
			//设置
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}
			nettool.SetYamuxConn(kcpconn)

			err = msg.WriteMsg(kcpconn, &models.Ping{})
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return nil, err
			}
			//:TODO 超时返回
			rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*2)
			if err != nil {
				kcpconn.Close()
				log.Println(err)
				return nil, err
			}
			switch m := rawMsg.(type) {
			case *models.Pong:
				{
					log.Println("get pong from p2p kcpconn")
					_ = m
					//TODO:认证
					config := yamux.DefaultConfig()
					//config.EnableKeepAlive = false
					p2pSubSession, err := yamux.Client(kcpconn, config)
					if err != nil {
						kcpconn.Close()
						log.Println("create sub session err:" + err.Error())
						return nil, err
					}
					return p2pSubSession, err
				}
			default:
				log.Println("type err")
			}
		}
	default:
		log.Println("type err")
		return nil, err
	}
	return nil, err
}
