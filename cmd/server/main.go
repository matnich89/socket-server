package main

import (
	"log"
	"net/http"
	"os"

	"syscallserver/internal/netsocket"
	"syscallserver/internal/service/proxy"
)

func main() {

	proxyTargetFlag := os.Getenv("TARGET_URL")

	if proxyTargetFlag == "" {
		log.Fatalln("target url not set")
	}

	reverseProxy, err := proxy.New(proxyTargetFlag)

	if err != nil {
		log.Fatalln(err)
	}

	server, err := netsocket.New("127.0.0.1", 8080)
	if err != nil {
		log.Fatalln(err)
	}

	go server.Listen()
	log.Println("http server started....")

	http.HandleFunc("/", reverseProxy.HandleProxyRequest)
	log.Println("reverse proxy starting....")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}
