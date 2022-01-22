# Arpspoofer

## ARP (Address Resolution Protocol)

ARP is a protocol that translates IP addresses into MAC addresses by broadcasting ARP requests with the IP address of the target onto the local network. If a device knows the corresponding MAC address of that target, it will send an ARP reply. Mappings between IP addresses and MAC addresses are stored locally on each device in the so called ARP table.

## ARP Spoofing

ARP Spoofing is a Man in the Middle attack that allows an attacker to intercept communications of a target device. ARP is a stateless protocol, i.e., it will accept ARP replies without having issued a corresponding ARP request beforehand. Furthermore, the MAC address that the target learns through the ARP reply does not have to be the correct one for the given IP address. This allows an attacker to set false MAC addresses in the ARP table for any device on the network.

<img src="https://www.thesslstore.com/blog/wp-content/uploads/2021/02/arp0.png">

For example, to intercept the communication between Bob and Mary, the attacker sends a spoofed ARP reply to Bob. It tells Bob that the MAC address of Mary (192.168.0.5) is the attacker's MAC. Then, the attacker can set the entry for Bob's IP (192.168.0.1) within Mary's ARP table to his MAC address using another spoofed ARP reply. As a result, when Bob and Mary want to commmunicate with each other, the packets will be sent to the attacker's machine, where they can be inspected and forwarded accordingly.

## Usage

```
sudo go run . -t <TARGET_IP> -r <ROUTER_IP> -i <INTERFACE>
```
