package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

var (
	wireguardIP             = flag.String("wgIP", "", "IP address for the Wireguard interface")
	wireguardPort           = flag.String("wgPort", "", "Port for the Wireguard interface")
	wireguardPrivateKey     = flag.String("wgPrivateKey", "", "Private key for the Wireguard interface")
	wireguardPeerPublicKey  = flag.String("wgPeerPublicKey", "", "Public key of the peer for the Wireguard tunnel")
	wireguardPeerEndpoint   = flag.String("wgPeerEndpoint", "", "Endpoint (IP and port) of the peer for the Wireguard tunnel")
	wireguardPeerAllowedIPs = flag.String("wgPeerAllowedIPs", "", "Allowed IPs for the peer for the Wireguard tunnel")

	nginxListenIP   = flag.String("nginxIP", "", "IP address for the nginx to listen on")
	nginxListenPort = flag.String("nginxPort", "", "Port for the nginx to listen on")
	nginxServerName = flag.String("nginxServerName", "", "Server name for the nginx server")

	test = flag.String("test", "", "Server name for the nginx server")
)

func main() {
	flag.Parse()
	fmt.Println(*wireguardIP, *wireguardPort, *wireguardPrivateKey, *wireguardPeerPublicKey, *wireguardPeerEndpoint, *wireguardPeerAllowedIPs, *nginxListenIP, *nginxListenPort)
	// check if all the required flags are passed or not
	if *wireguardIP == "" || *wireguardPort == "" || *wireguardPrivateKey == "" || *wireguardPeerPublicKey == "" || *wireguardPeerEndpoint == "" || *wireguardPeerAllowedIPs == "" || *nginxListenIP == "" || *nginxListenPort == "" || *nginxServerName == "" {
		log.Fatal("All flags are not provided")
	}

	// Create Wireguard config file
	wireguardConfig := []byte(fmt.Sprintf("[Interface]\nAddress = %s/32\nListenPort = %s\nPrivateKey = %s\n\n[Peer]\nPublicKey = %s\nEndpoint = %s\nAllowedIPs = %s\n",
		*wireguardIP, *wireguardPort, *wireguardPrivateKey, *wireguardPeerPublicKey, *wireguardPeerEndpoint, *wireguardPeerAllowedIPs))
	err := ioutil.WriteFile("wg0.conf", wireguardConfig, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wireguard config file created")

	// Start Wireguard tunnel
	cmd := exec.Command("wg-quick", "up", "./wg0.conf")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wireguard tunnel established")

	// Create nginx config file
	nginxConfig := []byte(fmt.Sprintf("events {\n    worker_connections 1024;\n}\n\nhttp {\n    server {\n        listen %s:%s;\n        server_name %s;\n\n        location / {\n            proxy_pass http://%s:%s;\n        }\n    }\n}",
		*nginxListenIP, *nginxListenPort, *nginxServerName, *wireguardIP, *wireguardPort))
	err = ioutil.WriteFile("nginx.conf", nginxConfig, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nginx config file created")

	// Start nginx proxy
	cmd = exec.Command("nginx", "-c", "/nginx.conf")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nginx proxy started")

	// cleanup the files
	os.Remove("wg0.conf")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("Press 'CTRL+C' to exit...")
	<-sig

	fmt.Println("Exiting...")
}
