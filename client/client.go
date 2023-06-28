package client

import (
	"fmt"
	"net"
	C "rdcp/common"
	"strconv"
)

type Client struct {
	Host       string
	Port       int
	connection net.Conn
}

func NewClient(host string, port int) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}

func (c *Client) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.GetAddress())
	if err != nil {
		return err
	}

	c.connection = conn
	return nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}

func (c *Client) SendRequest(req *C.Request) (*C.Response, error) {
	_, err := c.connection.Write(req.GetSerialized())
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1024)
	_, e := c.connection.Read(b)
	if e != nil {
		return nil, e
	}

	return C.NewResponseFromString(string(b))
}

func sendArbitraryRequest(host string, m C.METHOD, args ...string) (*C.Response, error) {
	h, p, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}
	pi, e := strconv.Atoi(p)
	if e != nil {
		return nil, e
	}

	cl := NewClient(h, pi)
	err = cl.Connect()
	if err != nil {
		return nil, err
	}
	defer cl.Close()

	req := C.NewRequest(m)
	return cl.SendRequest(req)
}

func Order(host string, order string) (*C.Response, error) {
	return sendArbitraryRequest(host, C.ORDER, order)
}
