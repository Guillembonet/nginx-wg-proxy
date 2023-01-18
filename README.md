# nginx-wg-proxy

[![Docker Version](https://img.shields.io/docker/v/bunetz/nginx-wg-proxy?sort=date)](https://hub.docker.com/r/bunetz/nginx-wg-proxy)
[![Docker Pulls](https://img.shields.io/docker/pulls/bunetz/nginx-wg-proxy)](https://hub.docker.com/r/bunetz/nginx-wg-proxy)

### **This code has been almost fully generated by Chat GPT ([conversation here](https://sharegpt.com/c/wrLCMpx)).**
## Description
Small docker container which allows proxying http requests through a wireguard tunnel. Useful to expose your home computer local app with a public endpoint.

## Usage
The container needs to run with `NET_ADMIN` capability enabled, it doesn't need to run with host networking.

1. Once the container is running check logs to copy the wireguard configuration to your peer. Copy the contents of:
```
*** Wireguard config for peer ***
<contents>
*********      End      *********
```
2. Paste it in a file named `wg0.conf`
3. Run `sudo wg-quick up ./wg0.conf` in the file location
4. You may need to allow connections to your `nginxProxyPort` in the firewall.
### Flags:
You need to set IP and port for wgEndpoint or WgPeerEndpoint.

You also need to specify values for the flags which don't have defaults.
```
  -nginxIP string
        IP address for the nginx to listen on (default "0.0.0.0")
  -nginxPort string
        Port for the nginx to listen on (default "8080")
  -nginxProxyPort string
        Port for the nginx to proxy to (default "8080")
  -nginxServerName string
        Server name for the nginx server (default "wg-proxy")
  -wgEndpointIP string
        Peer endpoint IP used by the peer for the Wireguard tunnel
  -wgEndpointPort string
        Peer endpoint port used by the peer for the Wireguard tunnel (default "52122")
  -wgIP string
        IP address for the Wireguard interface (default "10.0.0.1")
  -wgPeerAllowedIPs string
        Allowed IPs for the peer for the Wireguard tunnel (default "10.0.0.2/32")
  -wgPeerEndpointIP string
        Peer endpoint IP used by the host for the Wireguard tunnel
  -wgPeerEndpointPort string
        Peer endpoint port used by the host for the Wireguard tunnel (default "52122")
  -wgPeerPublicKey string
        Public key of the peer for the Wireguard tunnel
  -wgPeerWireguardPort string
        Port for the Wireguard interface of the peer (default "52122")
  -wgPort string
        Port for the Wireguard interface (default "52122")
  -wgPrivateKey string
        Private key for the Wireguard interface
```
