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
	id := flag.Int("id", 0, "Identifier 1-65535")

	flag.Parse()

	helpFlag := flag.Bool("help", false, "Display help message")

	// Parse the command-line flags
	flag.Parse()

	if *id == 0 {
		printHelp()
		os.Exit(0)
	}

	// Help flag
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	// Destination IP address
	listenOn := *listen
	identifier := *id

	// Create a connection
	conn, err := icmp.ListenPacket("ip4:icmp", listenOn)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer conn.Close()

	fmt.Printf("Listening for ICMP Echo Request packets on %v with identifier %v\n", listenOn, identifier)

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
				// Check if the identifier matches the desired identifier
				if body.ID == identifier {
					name := string(body.Data[:100])
					// Delete the padding
					name = strings.TrimRight(name, "\x00")
					data := body.Data[100:]
					fmt.Println("Received ICMP Echo Request from:", src.String(), " filename: ", name, ", data size:", len(data))
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
