package server

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"os"
	C "rdcp/common"
	"time"
)

type RDCPServer struct {
	Host string
	Port string

	listener net.Listener
	Handlers map[string]func(*C.Request, *C.Response) error

	tlsConfig *tls.Config
}

func NewRDCPServer(host string, port string) *RDCPServer {
	return &RDCPServer{
		Host:     host,
		Port:     port,
		Handlers: make(map[string]func(*C.Request, *C.Response) error),
	}
}

func (s *RDCPServer) Listen() error {
	listener, err := net.Listen("tcp", s.GetAddress())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	s.listener = listener
	return s.listenIncomingConnection()
}

func (s *RDCPServer) listenIncomingConnection() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		go s.handleIncomingConnection(conn)
	}
}

func (s *RDCPServer) Handle(method C.METHOD, handler func(*C.Request, *C.Response) error) {
	s.Handlers[method] = handler
}

func (s *RDCPServer) WithTLS(tlsConfig *tls.Config) {
	go s.listenIncomingConnection()
}

func (s *RDCPServer) GetAddress() string {
	return s.Host + ":" + s.Port
}

func (s *RDCPServer) Close() {
	s.listener.Close()
}

func (s *RDCPServer) handleIncomingConnection(conn net.Conn) {

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		log.Printf("Received: %s", b[:n])

		response := &bytes.Buffer{}

		r := C.ParseContentIntoRequest(b)
		w := &C.Response{
			Headers:   make(map[string]string),
			Method:    "R/" + r.Method,
			Timestamp: int64(time.Now().Unix()),
			Body:      response,
		}

		// call the handler for that method
		if handler, ok := s.Handlers[r.Method]; ok {
			handler(r, w)
			conn.Write(response.Bytes())
		}
	}
}
