package main

import (
	"fmt"
	C "rdcp/common"
	S "rdcp/server"
)

func main() {
	server := S.NewRDCPServer("127.0.0.1", "8080")

	server.HandleOrder("Action1", func(req *C.Request, res *C.Response) error {
		p := req.GetActionParameters()
		fmt.Printf("p: %v\n", p)
		res.Write([]byte("Hello World"))
		return nil
	})

	fmt.Println("Listening")
	server.Listen()
}
