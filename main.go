package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/shangsuru/arpspoof/arpspoof"
)

func main() {
	argTargetIP := flag.String("t", "", "The IP address of the target")
	argRouterIP := flag.String("r", "", "The IP address of the router")
	networkInterface := flag.String("i", "wlan0", "The network interface")
	flag.Parse()

	targetIP := net.ParseIP(*argTargetIP)
	routerIP := net.ParseIP(*argRouterIP)

	if targetIP == nil || routerIP == nil {
		fmt.Println("Usage:  -t <TARGET_IP> -r <ROUTER_IP> -i <INTERFACE>")
		os.Exit(0)
	}

	arpspoof.ArpSpoof(*networkInterface, routerIP, targetIP)
}
