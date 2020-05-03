package session

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

func NewP2PCtrlAsServer(newstream net.Conn, TokenModel *models.TokenClaims) (*yamux.Session, error) {
	err := msg.WriteMsg(newstream, &models.ReqNewP2PCtrlAsClient{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rawMsg, err := msg.ReadMsgWithTimeOut(newstream, time.Second*5)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	switch m := rawMsg.(type) {
	case *models.ReqNewP2PCtrlAsServer:
		{
			//监听一个随机端口号，接受P2P方的连接
			externalUDPAddr, listener, err := p2p.GetP2PListener(TokenModel)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			p2p.SendPackToPeerByReqNewP2PCtrlAsServer(listener, m)
			//开始转kcp监听
			go func() {
				time.Sleep(time.Second)
				//TODO：发送认证码用于后续校验
				msg.WriteMsg(newstream, &models.RemoteNetInfo{
					IntranetIp:   listener.LocalAddr().(*net.UDPAddr).IP.String(),
					IntranetPort: listener.LocalAddr().(*net.UDPAddr).Port,
					ExternalIp:   externalUDPAddr.IP.String(),
					ExternalPort: externalUDPAddr.Port,
				})
				//TODO:这里控制连接的处理？
				newstream.Close()
			}()
			return kcpListener(listener)
		}
	default:
		log.Println("不是ReqNewP2PCtrlAsServer")
		return nil, errors.New("不是ReqNewP2PCtrlAsServer")
	}
}

//TODO：listener转kcp服务侦听
func kcpListener(listener *net.UDPConn) (*yamux.Session, error) {
	kcplis, err := kcp.ServeConn(nil, 10, 3, listener)
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
		kcplis.Close()
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
	//b:=make([]byte,1024)
	//n,err:=conn.Read(b)
	//log.Println(string(b[0:n]))
	//lis.Close()
	//	从从conn中读取p2p另一方发来的认证消息，认证成功之后包装为mux服务端
	err = kcplis.SetDeadline(time.Time{})
	if err != nil {
		kcplis.Close()
	}
	return kcpConnHdl(kcpconn)
}

func kcpConnHdl(kcpconn net.Conn) (*yamux.Session, error) {
	//:TODO 超时返回
	rawMsg, err := msg.ReadMsgWithTimeOut(kcpconn, time.Second*2)
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
			err := msg.WriteMsg(kcpconn, &models.Pong{})
			if err != nil {
				kcpconn.Close()
				log.Println(err.Error())
				return nil, err
			}
			config := yamux.DefaultConfig()
			config.EnableKeepAlive = false
			session, err := yamux.Client(kcpconn, config)
			if err != nil {
				kcpconn.Close()
				log.Println(err.Error())
				return nil, err
			}
			return session, err
		}
	default:
		log.Println("获取到了一个未知的P2P握手消息")
	}
	return nil, err
}
