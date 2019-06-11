
package main

import (
	"flag"
	"fmt"
	"shorty-url/db"
	"shorty-url/server"

	"github.com/gorilla/mux"
)

var port = flag.Int("port", 8080, "port to host server on")

func serve(a server.App) {
	a.Router = mux.NewRouter()
	a.Init()
	a.Run(fmt.Sprintf(":%d",*port))
}

func main() {
	flag.Parse()
	a := server.App{}
	a.DB = db.GetInMemDb()
	serve(a)
}