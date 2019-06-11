package main

import (
	"flag"
	"fmt"
	"shorty-url/db"
	"shorty-url/server"
)

var port = flag.Int("port", 8080, "port to host server on")

func serve(a server.ShortyServer) {
	a.Init()
	a.Run(fmt.Sprintf(":%d", *port))
}

func main() {
	flag.Parse()
	a := server.CreateServer(db.GetInMemDb())
	serve(a)
}
