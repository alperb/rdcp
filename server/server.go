package server

import (
	"crypto/tls"
	"fmt"
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
	Handlers map[string]map[string]func(*C.Request, *C.Response) error

	Options map[string]interface{}

	tlsConfig *tls.Config
	Version   string
}

func NewRDCPServer(host string, port string) *RDCPServer {
	s := &RDCPServer{
		Host:     host,
		Port:     port,
		Handlers: make(map[string]map[string]func(*C.Request, *C.Response) error),
		Options:  make(map[string]interface{}),
		Version:  C.BLIZZARD,
	}

	s.Handlers[C.ORDER] = make(map[string]func(*C.Request, *C.Response) error)
	s.Handlers[C.OPTIONS] = make(map[string]func(*C.Request, *C.Response) error)

	return s
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

func (s *RDCPServer) HandleOrder(action string, handler func(*C.Request, *C.Response) error) {
	s.Handlers[C.ORDER][action] = handler
}

func (s *RDCPServer) HandleOptions(action string, handler func(*C.Request, *C.Response) error) {
	s.Handlers[C.OPTIONS][action] = handler
}

func (s *RDCPServer) WithTLS(tlsConfig *tls.Config) *RDCPServer {
	// go s.listenIncomingConnection()
	return s
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
			conn.Write([]byte("Unexpected behavior"))
			conn.Close()
		}

		log.Printf("Received: %s", b[:n])

		response := C.RequestBody{}

		r := C.ParseContentIntoRequest(b)

		// if the method is not supported, return an error
		if r.Method == "" {
			conn.Write([]byte("Method not supported"))
			conn.Close()
			return
		}

		// check for method
		switch r.Method {
		case C.HEARTBEAT:
			res := C.NewResponse(C.HEARTBEAT)
			res.AddHeader("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
			res.AddHeader("From", conn.RemoteAddr().String())
			res.Write([]byte("OK"))

			conn.Write(res.GetSerialized())
		case C.ORDER:
			// create a response object for handler
			w := &C.Response{
				Headers:   C.CreateDefaultHeaders(time.Now(), s.Version, true),
				Method:    "R/" + r.Method,
				Timestamp: int64(time.Now().Unix()),
				Body:      response,
			}

			// call the handler for that method
			if handler, ok := s.Handlers[r.Method][r.GetActionName()]; ok {
				handler(r, w)
				conn.Write(w.GetSerialized())
			} else {
				conn.Write([]byte("Action not handled\r\n"))
			}
		}

	}
}
