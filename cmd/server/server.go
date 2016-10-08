package main

import (
	"net"

	"github.com/Sirupsen/logrus"

	"github.com/thetooth/curveface"
	"github.com/thetooth/curveface/etcd"
)

const bind = "0.0.0.0:9001"

func main() {
	log := logrus.New()

	// Connect to etcd
	c, err := etcd.Dial("http://127.0.0.1:4001")
	if err != nil {
		log.Fatal(err)
	}
	c.MkDir("/curveface")
	c.MkDir("/curveface/server")

	// Configure server
	serverPrivkey, err := curveface.PrivkeyFromEtcd(c, "/curveface/server/privkey")
	if err != nil {
		log.Fatalf("Failed to parse server private key: %s", err)
	}
	serverPubkey, err := curveface.PubkeyFromEtcd(c, "/curveface/server/pubkey")
	if err != nil {
		log.Fatalf("Failed to parse server public key: %s", err)
	}

	// Load client certs
	keys, err := curveface.KeysFromEtcd(c)
	if err != nil {
		log.Fatalf("Failed to parse list of client keys: %s", err)
	}

	log.Infof("%d keys in-memory, listening for connections...", len(keys))

	// Wait for connections
	for listener, _ := net.Listen("tcp4", bind); listener != nil; {
		// Get secure socket
		ssocket, err := curveface.Accept(listener, serverPrivkey, serverPubkey, keys)
		if err != nil {
			log.Warnf("Failed to process authorization: %s", err)
			continue
		}

		var packet [24]byte

		// Receive client message
		_, err = ssocket.Read(packet[:])
		if err != nil {
			log.Warnf("Failed to read from wrapped socket: %s", err)
			ssocket.Close()
			continue
		}

		// Handle
		switch {
		case string(packet[:5]) == "gface":
			_, err = ssocket.Write([]byte("( ≖‿≖)"))
			if err != nil {
				log.Warnf("Failed to write to wrapped socket: %s", err)
			}
		}

		// Close on exit
		if err := ssocket.Close(); err != nil {
			log.Warnf("Failed to close socket: %s", err)
		}

	}
}
