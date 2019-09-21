package rtsp

import (
	"rtspServer/setting"
	"rtspServer/tlswrap"
	"fmt"
	"log"
	"net"
	"sync"

	// "rtspServer/penggy/EasyGoLib/utils"
)

type Server struct {
	// TCPListener *net.TCPListener
	TCPListener net.Listener
	// 是否使用tls
	Tls         bool
	TCPPort     int
	Stoped      bool
	pushers     map[string]*Pusher // Path <-> Pusher
	pushersLock sync.RWMutex
}

var Instance *Server = &Server{
	Stoped:  true,
	// 默认不使用tls
	Tls: false,
	// 可执行文件名必须为easydarwin(不区分大小写),配置文件为可执行文件名+".ini"(已固定为easydarwin.ini)
	TCPPort: setting.Conf().Section("rtsp").Key("port").MustInt(554),
	pushers: make(map[string]*Pusher),
}

func GetServer() *Server {
	return Instance
}

func NewServer(port int, tls bool) *Server {
	return &Server{
		Stoped: true,
		Tls: tls,
		TCPPort: port,
		pushers: make(map[string]*Pusher),
	}
}


func (server *Server) Start() (err error) {

	// 配置文件内有tls时,启用rtspsServer; 没有时rtspsServer=nil
	err = nil
	if server == nil {
		return nil
	}

	port := server.TCPPort
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	var listener net.Listener
	// var listener *net.TCPListener
	if server.Tls {
		cert := setting.Conf().Section("tls").Key("cert").MustString("")
		key := setting.Conf().Section("tls").Key("key").MustString("")
		listener, err = tlswrap.NewTlsListener(cert, key, addr)
		// listener, err = tlswrap.NewTlsListener2(cert, key, port)
		log.Println("rtsps server start on", server.TCPPort)
	} else {
		listener, err = net.ListenTCP("tcp", addr)
		log.Println("rtsp server start on", server.TCPPort)
	}

	if err != nil {
		return
	}

	server.Stoped = false
	server.TCPListener = listener
	// log.Println("rtsp server start on", server.TCPPort)
	// networkBuffer := setting.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	for !server.Stoped {
		// conn, err := server.TCPListener.AcceptTCP()
		conn, err := server.TCPListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// if err := conn.SetReadBuffer(networkBuffer); err != nil {
		// 	log.Printf("rtsp server conn set read buffer error, %v", err)
		// }
		// if err := conn.SetWriteBuffer(networkBuffer); err != nil {
		// 	log.Printf("rtsp server conn set write buffer error, %v", err)
		// }
		session := NewSession(server, conn)
		go session.Start()
	}
	return
}

func (server *Server) Stop() {
	log.Println("rtsp server stop on", server.TCPPort)
	server.Stoped = true
	if server.TCPListener != nil {
		server.TCPListener.Close()
		server.TCPListener = nil
	}
	server.pushersLock.Lock()
	server.pushers = make(map[string]*Pusher)
	server.pushersLock.Unlock()
}

func (server *Server) AddPusher(pusher *Pusher) {
	server.pushersLock.Lock()
	if _, ok := server.pushers[pusher.Path]; !ok {
		server.pushers[pusher.Path] = pusher
		go pusher.Start()
		log.Printf("%v start, now pusher size[%d]", pusher, len(server.pushers))
	}
	server.pushersLock.Unlock()
}

func (server *Server) RemovePusher(pusher *Pusher) {
	server.pushersLock.Lock()
	if _pusher, ok := server.pushers[pusher.Path]; ok && pusher.ID == _pusher.ID {
		delete(server.pushers, pusher.Path)
		log.Printf("%v end, now pusher size[%d]\n", pusher, len(server.pushers))
	}
	server.pushersLock.Unlock()
}

func (server *Server) GetPusher(path string) (pusher *Pusher) {
	server.pushersLock.RLock()
	pusher = server.pushers[path]
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPushers() (pushers map[string]*Pusher) {
	pushers = make(map[string]*Pusher)
	server.pushersLock.RLock()
	for k, v := range server.pushers {
		pushers[k] = v
	}
	server.pushersLock.RUnlock()
	return
}

func (server *Server) GetPusherSize() (size int) {
	server.pushersLock.RLock()
	size = len(server.pushers)
	server.pushersLock.RUnlock()
	return
}
