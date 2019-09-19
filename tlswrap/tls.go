package tlswrap

import (
	// "bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)
/*
func TlsServer() {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
*/

// func NewTlsListener(certFile string, keyFile string, addr *net.TCPAddr) (*net.TCPListener, error) {
func NewTlsListener(certFile string, keyFile string, addr *net.TCPAddr) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// *TCPListener
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ln := tls.NewListener(listener, config)
	return ln, nil
}

// func NewTlsListener2(certFile string, keyFile string, port int) (*net.TCPListener, error) {
func NewTlsListener2(certFile string, keyFile string, port int) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	laddr := fmt.Sprintf(":%d", port)
	ln, err := tls.Listen("tcp", laddr, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ln, nil

	// listener, ok := ln.(*net.TCPListener)
	// if ok {
	// 	// false
	// 	log.Printf("类型转换:%v\n", ok)
	// }
	// return listener, nil
}
