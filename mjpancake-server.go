package main

import (
	"os"

	"github.com/mjpancake/mjpancake-server/srv"
)

func main() {
	port := "6171"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	srv.Serve(port)
}
