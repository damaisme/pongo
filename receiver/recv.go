package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	listen := flag.String("l", "127.0.0.1", "Listen IP address")
	p := flag.String("p", "", "Password for identifier")
	helpFlag := flag.Bool("help", false, "Display help message")

	flag.Parse()

	// Parse the command-line flags
	flag.Parse()

	if *p == "" {
		printHelp()
		os.Exit(0)
	}

	// Help flag
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	password := *p

	// Destination IP address
	listenOn := *listen

	// Create a connection
	conn, err := icmp.ListenPacket("ip4:icmp", listenOn)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer conn.Close()

	fmt.Printf("Listening for ICMP Echo Request packets on %v with identifier %v\n", listenOn, password)

	// Buffer to read incoming ICMP packets
	buffer := make([]byte, 60000)

	for {
		// Read incoming ICMP packet
		n, src, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal("Error reading ICMP packet:", err)
		}

		// Parse the ICMP packet
		packet, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), buffer[:n])
		if err != nil {
			log.Fatal("Error parsing ICMP packet:", err)
		}

		// Check if the received packet is an ICMP Echo Request
		if packet.Type == ipv4.ICMPTypeEcho && packet.Code == 0 {
			// Retrieve the body of the ICMP message
			switch body := packet.Body.(type) {
			case *icmp.Echo:

				icmpPass := string(body.Data[:50])
				name := string(body.Data[50:150])
				// Delete the padding
				name = strings.TrimRight(name, "\x00")
				icmpPass = strings.TrimRight(icmpPass, "\x00")
				// Check if the password matches
				if icmpPass == password {
					// Get the chunk
					data := body.Data[150:]
					fmt.Println("Received ICMP Echo Request from:", src.String(), ", password:", icmpPass, ", filename:", name, ", data size:", len(data))

					receivePath := "received_" + name
					// Open the file in append mode, create it if it doesn't exist
					file, err := os.OpenFile(string(receivePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						log.Fatal("Error opening file:", err)
					}
					defer file.Close()

					fmt.Println("Writing to file:", receivePath)
					// Append to the file
					_, err = file.Write(data)
					if err != nil {
						log.Fatal("Error appending to file:", err)
					}

				}
			}
		}
	}
}

func printHelp() {
	fmt.Println("Receiving file through ICMP")
	fmt.Println("Usage: recv.go [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
