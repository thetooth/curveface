# Curveface
![Image of currency](https://thetooth.name/images/gface-royal-cropped.png)
A library and server for receiving signed gface payloads via SPKI. Uses elliptic curve encryption of communications between peers with [ZeroMQ ZMTP](https://rfc.zeromq.org/spec:37/ZMTP/) framing for reliable gface transactions.

## Features
 * [Robust encryption](https://godoc.org/golang.org/x/crypto/nacl/box)
 * Distributed and highly available by design
 * Easy to use, one-liner for gface requests:
<br>`GetGface(endpoint, priv, pub, srvpub) ([12]byte, error)`

## Roadmap
- [ ] Automatic expiry of in-memory keys
- [x] IRC bot
- [ ] C/C++ library

## Building
The library is go gettable. To build the server and included tools run the following:
```shell
$ cd $GOPATH/src/github.com/thetooth/curveface/
$ go build -v cmd/genkeypair/genkeypair.go
$ go build -v cmd/server/server.go
```
You can then use the genkeypair tool to create a set of keys for the server and clients.