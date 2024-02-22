package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host   = flag.String("host", "0.0.0.0", "Host to serve on")
	port   = flag.Int("port", 8080, "Port to serve on")
	folder = flag.String("dir", ".", "Folder to serve")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", logRequest(handleRequest))
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Serving files from %s on %s\n", *folder, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
