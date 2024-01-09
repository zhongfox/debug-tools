package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var httpPorts = os.Getenv("HTTP_PORTS")
var udpPorts = os.Getenv("UDP_PORTS")
var tcpPorts = os.Getenv("TCP_PORTS")
var service = os.Getenv("SERVICE")
var podName = os.Getenv("POD_NAME")
var podIP = os.Getenv("POD_IP")

func main() {
	if service == "" {
		service = "unknown"
	}

	startHTTPServers()
	startUDPServers()
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
			//mux.HandleFunc(fmt.Sprintf("/%s", service), func(rw http.ResponseWriter, req *http.Request) {
			mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				log.Printf("=>Receiving call (port:%s, path:%s) from %s\n", port, req.URL.Path, req.RemoteAddr)
				log.Println("  headers:")
				for name, values := range req.Header {
					for _, value := range values {
						log.Printf("    %s: %s \n", name, value)
					}
				}

				rw.Write([]byte(fmt.Sprintf("%s(%s %s:%s) received from %s on path %s\n", service, podName, podIP, port, req.RemoteAddr, req.URL.Path)))
			})
			server := http.Server{Addr: ":" + port, Handler: mux}
			log.Printf("starting http service %s on port %s \n", service, port)
			log.Fatal(server.ListenAndServe())
		}()
	}

}

func startUDPServers() {
	ports := udpPorts
	if ports == "" {
		ports = "4100"
	}

	ps := strings.Split(ports, ",")

	for _, p := range ps {
		port, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		go func() {
			log.Printf("starting udp service %s on port %d \n", service, port)
			l, err := net.ListenUDP("udp4", &net.UDPAddr{
				Port: port,
				IP:   net.ParseIP("0.0.0.0"),
			})
			if err != nil {
				log.Fatal(err)
			}
			defer l.Close()

			// Read from UDP listener in endless loop
			for {
				var buf [512]byte
				_, addr, err := l.ReadFromUDP(buf[0:])
				if err != nil {
					log.Printf("udp read err: %v\n", err)
					return
				}

				fmt.Print("> ", string(buf[0:]))

				// Write back the message over UPD
				l.WriteToUDP([]byte("Hello UDP Client\n"), addr)
			}

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
