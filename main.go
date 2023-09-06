package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"io/ioutil"
	"github.com/akamensky/argparse"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main (){
  parser := argparse.NewParser("pongo", "Transfer file over ICMP")
  commandSend := parser.NewCommand("send", "Send a file using ICMP packets")
  commandRecv := parser.NewCommand("recv", "Receive a file sent via ICMP packets.")


  dst := commandSend.String(
    "d", "dst", &argparse.Options{
        Required: true,
        Help: "Destination host where the file will be received.",
    },
  )

  file := commandSend.String(
    "f", "file", &argparse.Options{
        Required: true,
        Help: "Path to the file you want to send.",
    },
  )

  mpacket := commandSend.Int(
    "m", "msize", &argparse.Options{
        Help: "Maximum packet size per ICMP Message in bytes (default is 8980).",
        Default: 8980,
    },
  )

  rsecret := commandSend.String(
    "s", "secret", &argparse.Options{
        Required: true,
        Help: "Identifier shared between the sender and receiver.",
    },
  )

	listen := commandRecv.String(
    "l", "listen", &argparse.Options{
        Required: true,
        Help: "IP address and port to listen for incoming ICMP packets.",
    },
  )

  ssecret := commandRecv.String(
    "s", "secret", &argparse.Options{
        Required: true,
        Help: "Identifier shared between the sender and receiver.",
    },
  )

	opath := commandRecv.String(
    "o", "opath", &argparse.Options{
    		Default: ".",
        Help: "Output path.",
    },
  )


  err := parser.Parse(os.Args)
  if err != nil {
      log.Fatalln(parser.Usage(err))
      return
  }

  switch {
  case commandSend.Happened():
      // fmt.Println("Send:", commandSend.Happened())
      send(*dst, *file, * mpacket, *ssecret)
  case commandRecv.Happened():
      // fmt.Println("Recv:", commandRecv.Happened())
			recv(*rsecret, *listen, *opath)
	}
}

func send(dst string, file string, mpacket int, ssecret string,) {

	// Create a connection
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal("Error creating connection:", err)
		return
	}
	defer conn.Close()

	// Read the file
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Split the file into smaller chunks
	chunks := splitData(fileData, mpacket-150)

	// Send each chunk to the destination
	for _, chunk := range chunks {

		// get file name
		parts := strings.Split(file, "/")
		fileName := parts[len(parts)-1]

		nameBytes := []byte(fileName)
		ssecretBytes := []byte(ssecret)

		// Check the length of the ssecret
		if len(ssecretBytes) > 50 {
			ssecretBytes = ssecretBytes[:50]
		} else if len(ssecretBytes) < 50 {
			// If the name is shorter than 50 bytes, add padding
			padding := make([]byte, 50-len(ssecretBytes))
			ssecretBytes = append(ssecretBytes, padding...)
		}

		// Check the length of the name
		if len(nameBytes) > 100 {
			nameBytes = nameBytes[:100]
		} else if len(nameBytes) < 100 {
			// If the name is shorter than 100 bytes, add padding
			padding := make([]byte, 100-len(nameBytes))
			nameBytes = append(nameBytes, padding...)
		}

		// Create the final message
		bodyData := append(ssecretBytes, nameBytes...)
		bodyData = append(bodyData, chunk...)

		// Create ICMP message
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   12345, // Set the ICMP identifier value here
				Seq:  1,
				Data: bodyData,
			},
		}

		// Serialize the ICMP message
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			log.Fatal("Error marshaling ICMP message:", err)
			return
		}

		// Resolve destination IP address
		destIP := net.ParseIP(dst)
		if destIP == nil {
			log.Fatal("Invalid destination IP address")
			return
		}

		// Send ICMP message to destination
		_, err = conn.WriteTo(msgBytes, &net.IPAddr{IP: destIP})
		if err != nil {
			log.Fatal("Error sending ICMP message:", err)
			return
		}

		fmt.Printf("%v bytes ICMP packet sent successfully to %v!\n", len(chunk), destIP)

		time.Sleep(200 * time.Millisecond)
	}
}

// Split data into smaller chunks
func splitData(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func recv(rsecret string, listen string, opath string){
	// Create a connection
	conn, err := icmp.ListenPacket("ip4:icmp", listen)
	if err != nil {
		log.Fatal("Error creating connection:", err)
	}
	defer conn.Close()

	fmt.Printf("Listening for ICMP Echo Request packets on %v with identifier %v\n", listen, rsecret)

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
				if len(body.Data) < 150 {
					continue
				}
				icmpPass := string(body.Data[:50])
				name := string(body.Data[50:150])
				// Delete the padding
				name = strings.TrimRight(name, "\x00")
				icmpPass = strings.TrimRight(icmpPass, "\x00")
				// Check if the secret matches
				if icmpPass == rsecret {
					// Get the chunk
					data := body.Data[150:]
					fmt.Println("Received ICMP Echo Request from:", src.String(), ", secret:", icmpPass, ", filename:", name, ", data size:", len(data))

					receivePath := opath + name
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

	fmt.Println(rsecret)
}
