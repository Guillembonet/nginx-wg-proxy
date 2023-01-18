package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	wireguardIP                = flag.String("wgIP", "10.0.0.1", "IP address for the Wireguard interface")
	wireguardPort              = flag.String("wgPort", "52122", "Port for the Wireguard interface")
	wireguardPrivateKey        = flag.String("wgPrivateKey", "", "Private key for the Wireguard interface")
	wireguardEndpointIP        = flag.String("wgEndpointIP", "", "Peer endpoint IP used by the peer for the Wireguard tunnel")
	wireguardEndpointPort      = flag.String("wgEndpointPort", "52122", "Peer endpoint port used by the peer for the Wireguard tunnel")
	wireguardPeerPublicKey     = flag.String("wgPeerPublicKey", "", "Public key of the peer for the Wireguard tunnel")
	wireguardPeerWireguardPort = flag.String("wgPeerWireguardPort", "52122", "Port for the Wireguard interface of the peer")
	wireguardPeerEndpointIP    = flag.String("wgPeerEndpointIP", "", "Peer endpoint IP used by the host for the Wireguard tunnel")
	wireguardPeerEndpointPort  = flag.String("wgPeerEndpointPort", "52122", "Peer endpoint port used by the host for the Wireguard tunnel")
	wireguardPeerAllowedIPs    = flag.String("wgPeerAllowedIPs", "10.0.0.2/32", "Allowed IPs for the peer for the Wireguard tunnel")

	nginxListenIP   = flag.String("nginxIP", "0.0.0.0", "IP address for the nginx to listen on")
	nginxListenPort = flag.String("nginxPort", "8080", "Port for the nginx to listen on")
	nginxServerName = flag.String("nginxServerName", "wg-proxy", "Server name for the nginx server")
	nginxProxyPort  = flag.String("nginxProxyPort", "8080", "Port for the nginx to proxy to")
)

func main() {
	flag.Parse()
	// check if all the required flags are passed or not
	if *wireguardIP == "" || *wireguardPort == "" || *wireguardPrivateKey == "" ||
		*wireguardPeerPublicKey == "" || *wireguardPeerWireguardPort == "" || *wireguardPeerAllowedIPs == "" ||
		*nginxListenIP == "" || *nginxListenPort == "" ||
		*nginxServerName == "" || *nginxProxyPort == "" {
		log.Fatal("All flags are not provided")
	}

	if (*wireguardEndpointIP == "" || *wireguardEndpointPort == "") && (*wireguardPeerEndpointIP == "" || *wireguardPeerEndpointPort == "") {
		log.Fatal("Need to specify at least one endpoint ip and port")
	}

	// Create Wireguard config file
	wireguardConfig := fmt.Sprintf("[Interface]\nAddress = %s/32\nListenPort = %s\nPrivateKey = %s\n\n[Peer]\nPublicKey = %s\nAllowedIPs = %s\n",
		*wireguardIP, *wireguardPort, *wireguardPrivateKey, *wireguardPeerPublicKey, *wireguardPeerAllowedIPs)
	if *wireguardPeerEndpointIP != "" && *wireguardPeerEndpointPort != "" {
		wireguardConfig = fmt.Sprintf("%sEndpoint = %s:%s\nPersistentKeepalive = 25\n", wireguardConfig, *wireguardPeerEndpointIP, *wireguardPeerEndpointPort)
	}
	err := ioutil.WriteFile("wg0.conf", []byte(wireguardConfig), 0644)
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

	peerIpSplit := strings.Split(*wireguardPeerAllowedIPs, "/")
	if len(peerIpSplit) != 2 {
		log.Fatalf("peer allowed ips has bad format: %s", *wireguardPeerAllowedIPs)
	}
	// Create nginx config file
	nginxConfig := []byte(fmt.Sprintf("events {\n    worker_connections 1024;\n}\n\nhttp {\n    server {\n        listen %s:%s;\n        server_name %s;\n\n        location / {\n            proxy_pass http://%s:%s;\n        }\n    }\n}",
		*nginxListenIP, *nginxListenPort, *nginxServerName, peerIpSplit[0], *nginxProxyPort))
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

	key, err := wgtypes.ParseKey(*wireguardPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	// Print Wireguard config file for the peer
	fmt.Println("*** Wireguard config for peer ***")
	fmt.Printf("[Interface]\nAddress = %s\nListenPort = %s\nPrivateKey = privateKey\n\n[Peer]\nPublicKey = %s\nAllowedIPs = %s/32\n", *wireguardPeerAllowedIPs, *wireguardPeerWireguardPort, key.PublicKey().String(), *wireguardIP)
	if *wireguardEndpointIP != "" && *wireguardEndpointPort != "" {
		fmt.Printf("Endpoint = %s:%s\nPersistentKeepalive = 25\n", *wireguardEndpointIP, *wireguardEndpointPort)
	}
	fmt.Println("*********      End      *********")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("Press 'CTRL+C' to exit...")
	<-sig

	fmt.Println("Exiting...")
}
