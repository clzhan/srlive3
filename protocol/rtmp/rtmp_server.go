package rtmp

import (
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"github.com/clzhan/srlive3/log"
)

var (
	handler ServerHandler = new(DefaultServerHandler)
)


type Server struct {
	Addr string //监听地址

	Handler ServerHandler

	ReadTimeout  time.Duration //读超时
	WriteTimeout time.Duration //写超时
	Lock         *sync.Mutex
}


func ListenAndServer(addr string) error {

	srv:= &Server{
		Addr:         addr,
		ReadTimeout:  time.Duration(time.Second * 30),
		WriteTimeout: time.Duration(time.Second * 30),
		Lock:         new(sync.Mutex),
	}

	return srv.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {

	addr := srv.Addr
	if addr == "" {
		addr = ":1935"
	}
	if addr == "" {
		addr = ":1935"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("rtmp server start listen at :", addr)

	for i := 0; i < runtime.NumCPU(); i++ {
		go srv.loop(listener)
	}

	return nil
}

func (srv *Server) loop(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		netConn, e := l.Accept()
		if e != nil {
			log.Error("Accept error :", e)
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Errorf("rtmp: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		log.Info("Accept pid :", os.Getpid())
		//srv.SetCon(grw)
		tempDelay = 0
		go serve(srv, netConn)
	}
}


func serve(srv *Server, con net.Conn) {
	log.Info("Accept", con.RemoteAddr(), "->", con.LocalAddr())
	con.(*net.TCPConn).SetNoDelay(true)
	conn := newconn(con, srv)
	if !handshake1(conn.buf) {
		conn.Close()
		return
	}
	log.Info("handshake", con.RemoteAddr(), "->", con.LocalAddr(), "ok")
	log.Debug("readMessage")
	msg, err := readMessage(conn)
	if err != nil {
		log.Error("NetConnecton read error", err)
		conn.Close()
		return
	}

	cmd, ok := msg.(*ConnectMessage)
	if !ok || cmd.Command != "connect" {
		log.Error("NetConnecton Received Invalid ConnectMessage ", msg)
		conn.Close()
		return
	}
	conn.app = getString(cmd.Object, "app")

	conn.objectEncoding = int(getNumber(cmd.Object, "objectEncoding"))
	log.Debug(cmd)
	log.Info(con.RemoteAddr(), "->", con.LocalAddr(), cmd, conn.app, conn.objectEncoding)
	err = sendAckWinsize(conn, 512<<10)
	if err != nil {
		log.Error("NetConnecton sendAckWinsize error", err)
		conn.Close()
		return
	}
	err = sendPeerBandwidth(conn, 512<<10)
	if err != nil {
		log.Error("NetConnecton sendPeerBandwidth error", err)
		conn.Close()
		return
	}
	err = sendStreamBegin(conn)
	if err != nil {
		log.Error("NetConnecton sendStreamBegin error", err)
		conn.Close()
		return
	}
	err = sendConnectSuccess(conn)
	if err != nil {
		log.Error("NetConnecton sendConnectSuccess error", err)
		conn.Close()
		return
	}
	conn.connected = true

	newNetStream(conn, handler, nil).readLoop()
}

func getNumber(obj interface{}, key string) float64 {
	if v, exist := obj.(Map)[key]; exist {
		return v.(float64)
	}
	return 0.0
}
