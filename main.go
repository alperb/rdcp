package main

import (
	"fmt"
	C "rdcp/common"
	S "rdcp/server"
)

func main() {
	server := S.NewRDCPServer("127.0.0.1", "8080")

	server.Handle(C.ORDER, func(req *C.Request, res *C.Response) error {
		res.Write([]byte("Hello World"))
		return nil
	})

	fmt.Println("Listening")
	server.Listen()
}
