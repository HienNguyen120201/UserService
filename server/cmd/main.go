package main

import (
	"flag"
	"log"
	"net/http"
	"userservice/server/internal"
)

var (
	debugListenAddr = flag.String("debug-listen-addr", "127.0.0.1:7801", "HTTP listen address.")
)

func main() {
	internal.ConnectDB()
	go internal.RunGRPC()
	// keep go routine goes on
	log.Fatal(http.ListenAndServe(*debugListenAddr, nil))
}
