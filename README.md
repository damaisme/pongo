# go-icmp
File Transfer over ICMP

## Here's how it works
![file-tf-through-icmp](https://github.com/radendi/go-icmp/assets/73756341/18b76660-b684-4521-b7ad-e288193f31d1)


Read more: https://blog.dama.zip/2023/06/file-transfer-over-ping-icmp-messages.html

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
  -p string
        Password for identifier (default "icmp-go")
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
  -l string
        Listen IP address (default "127.0.0.1")
  -p string
        Password for identifier
```

## Example
Sender:
```
go run send.go -f image.jpg -d 127.0.0.1 -p pass12345 -s 1500
```

Receiver:
```
go run recv.go -l 127.0.0.1 -p pass12345 
```
The file will be saved in the receiver as `received_image.jpg`
