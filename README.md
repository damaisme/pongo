# go-icmp
File Transfer over ICMP

## How this works
![file-tf-through-icmp](https://github.com/radendi/go-icmp/assets/73756341/48c819a9-0a73-413b-a7f3-9993c8a73997)
Read more:

## Usage
Show help using `go run send.go -help` or `go run recv.go -help`

send.go :
```
Sending file through ICMP
Usage: send.go [options]
Options:
  -d string
        DestinationIP address (default "127.0.0.1")
  -f string
        File Path to send
  -help
        Display help message
  -id int
        Identifier 1-65535
  -s int
        Max packet size (default 8980)
```
recv.go:
```
Receiving file through ICMP
Usage: recv.go [options]
Options:
  -help
        Display help message
  -id int
        Identifier 1-65535
  -l string
        Listen IP address (default "127.0.0.1")
```

## Example
Sender:
```
go run send.go -f image.jpg -d 127.0.0.1 -id 12345 -s 1500
```

Receiver:
```
go run recv.go -l 127.0.0.1 -id 12345 
```
The file will be saved in the receiver as `received_image.jpg`
