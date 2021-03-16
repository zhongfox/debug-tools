package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var httpPorts = os.Getenv("HTTP_PORTS")
var tcpPorts = os.Getenv("TCP_PORTS")
var service = os.Getenv("SERVICE")
var podName = os.Getenv("POD_NAME")
var podIP = os.Getenv("POD_IP")

func main() {
	if service == "" {
		panic("miss env service, exiting")
	}

	startHTTPServers()
	startTCPServers()

	select {}
}

func startHTTPServers() {
	ports := httpPorts
	if ports == "" {
		ports = "7000"
	}

	ps := strings.Split(ports, ",")

	for _, p := range ps {
		port := p
		go func() {
			http.NewServeMux()

			mux := http.NewServeMux()
			mux.HandleFunc(fmt.Sprintf("/%s", service), func(rw http.ResponseWriter, req *http.Request) {
				log.Printf("=>Receiving call on %s from %s\n", port, req.RemoteAddr)
				log.Println("  headers:")
				for name, values := range req.Header {
					for _, value := range values {
						log.Printf("    %s: %s \n", name, value)
					}
				}

				rw.Write([]byte(fmt.Sprintf("%s(%s %s:%s) received\n", service, podName, podIP, port)))
			})
			server := http.Server{Addr: ":" + port, Handler: mux}
			log.Printf("starting http service %s on port %s \n", service, port)
			log.Fatal(server.ListenAndServe())
		}()
	}

}

func startTCPServers() {
	ports := tcpPorts
	if ports == "" {
		ports = "4000"
	}

	ps := strings.Split(ports, ",")

	for _, p := range ps {
		port := p
		go func() {
			log.Printf("starting tcp service %s on port %s \n", service, port)
			l, err := net.Listen("tcp4", ":"+port)
			if err != nil {
				log.Fatal(err)
			}
			defer l.Close()

			for {
				c, err := l.Accept()
				if err != nil {
					log.Printf("accept err: %v\n", err)
					return
				}
				go handleTCPConnection(port, c)
			}
		}()
	}
}

func handleTCPConnection(port string, c net.Conn) {
	client := c.RemoteAddr().String()
	log.Printf("TCP connection on %s from %s\n", port, client)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("read err: %v from %s\n", err, client)
			return
		}

		line := strings.TrimSpace(netData)
		if line == "quit" {
			break
		}
		log.Printf("=>Received call on %s from %s %s\n", port, client, line)
		c.Write([]byte(fmt.Sprintf("%s(%s %s) received\n", service, podName, podIP)))
	}
	log.Printf("Closing the connection %s\n", client)
	c.Close()
}
