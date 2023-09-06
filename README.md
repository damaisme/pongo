# Pongo
**Pongo** is an innovative ICMP (Internet Control Message Protocol) file transfer tool designed to allow users to transfer files over a network using ICMP packets. It is a lightweight and efficient solution for transferring files in scenarios where traditional methods may not be suitable.

## Here's how it works
![file-tf-through-icmp](https://github.com/radendi/go-icmp/assets/73756341/18b76660-b684-4521-b7ad-e288193f31d1)


Read more: https://blog.dama.zip/2023/06/file-transfer-over-ping-icmp-messages.html

## How to Get Started
### Pre-Built Binary (linux amd64 only)

Download the Binary:

Run the following command to download the pre-built binary for your platform (Linux AMD64 in this example):

```
wget https://github.com/radendi/pongo/releases/download/v1.0.0/pongo-linux-amd64
```

Make the binary executable:
```
chmod +x pongo-linux-amd64
```

Move the binary to a directory in your system's PATH, so you can run it from anywhere:
```
sudo mv pongo-linux-amd64 /usr/local/bin/pongo
```

Start Using Pongo:
```
pongo -h
```

You can now use Pongo to transfer files using ICMP packets.

### Build Manually from Source Code

Clone the Repository:

Run the following command to clone the Pongo repository from GitHub:
```
git clone https://github.com/radendi/pongo.git
```

Build Pongo using the go build command:
```
cd pongo
go build -o pongo
```
Install Pongo:

Move the generated binary to a directory in your system's PATH, so you can run it from anywhere:
```
sudo mv pongo /usr/local/bin/pongo
```

Start Using Pongo:
```
pongo -h
```
You can now use Pongo to transfer files using ICMP packets.

## Example Usage
Sender:
```
pongo send -f image.jpg -d 127.0.0.1 -s secret12345 -m 5000
```

Receiver:
```
pongo recv -l 127.0.0.1 -s secret12345 -o received/
```
The file will be saved in the receiver at `received/image.jpg`
