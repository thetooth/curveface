package curveface

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/thetooth/curveface/etcd"

	"github.com/Rudd-O/curvetls"
)

// PKI stores PKI table
type PKI map[string]KeyConf

// KeyConf stores PKI field data
type KeyConf struct {
	Name   string
	Pubkey curvetls.Pubkey
}

// Client session type
type Client struct {
}

// Dial connects to a given endpoint and returns a secure socket interface
func (c Client) Dial(endpoint, priv, pub, srvpub string) (*curvetls.EncryptedConn, error) {
	clientPrivkey, err := curvetls.PrivkeyFromString(priv)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse client private key")
	}
	clientPubkey, err := curvetls.PubkeyFromString(pub)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse client public key")
	}
	serverPubkey, err := curvetls.PubkeyFromString(srvpub)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse server public key")
	}

	socket, err := net.Dial("tcp4", endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to socket")
	}

	longNonce, err := curvetls.NewLongNonce()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate nonce")
	}
	ssocket, err := curvetls.WrapClient(socket, clientPrivkey, clientPubkey, serverPubkey, longNonce)
	if err != nil {
		if curvetls.IsAuthenticationError(err) {
			return nil, errors.Wrap(err, "Server says unauthorized")
		}
		return nil, errors.Wrap(err, "Failed to wrap socket")
	}
	return ssocket, nil
}

// Accept incoming connections
func Accept(listener net.Listener, serverPrivkey curvetls.Privkey, serverPubkey curvetls.Pubkey, keys PKI) (*curvetls.EncryptedConn, error) {
	log := logrus.New()

	socket, err := listener.Accept()
	if err != nil {
		log.Warnf("Failed to accept socket: %s", err)
		return nil, errors.Wrap(err, "Failed to accept socket")
	}

	log.Printf("Request from %s", socket.RemoteAddr())

	longNonce, err := curvetls.NewLongNonce()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate nonce")
	}
	authorizer, clientpubkey, err := curvetls.WrapServer(socket, serverPrivkey, serverPubkey, longNonce)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to wrap socket")
	}

	log.Debugf("Client's public key is %s", clientpubkey)

	var ssocket *curvetls.EncryptedConn

	if _, ok := keys[clientpubkey.String()]; ok {
		log.Debugf("Client authorized")
		ssocket, err = authorizer.Allow()
	} else {
		err = authorizer.Deny()
	}

	if err != nil {
		ssocket.Close()
		return nil, errors.Wrap(err, "Failed to process authorization")
	}

	return ssocket, nil
}

// PrivkeyFromEtcd retrieves and decodes a private key from Etcd
func PrivkeyFromEtcd(c etcd.EtcdClient, key string) (curvetls.Privkey, error) {
	key, err := c.Get(key)
	if err != nil {
		return [32]byte{}, err
	}
	return curvetls.PrivkeyFromString(key)
}

// PubkeyFromEtcd retrieves and decodes a public key from Etcd
func PubkeyFromEtcd(c etcd.EtcdClient, key string) (curvetls.Pubkey, error) {
	key, err := c.Get(key)
	if err != nil {
		return [32]byte{}, err
	}
	return curvetls.PubkeyFromString(key)
}

// KeysFromEtcd loads the PKI table from Etcd
func KeysFromEtcd(c etcd.EtcdClient) (PKI, error) {
	keys := PKI{}
	ls, err := c.Ls("/curveface/client")
	if err != nil {
		return nil, errors.Wrap(err, "Could not list client certs")
	}

	for _, key := range ls {
		name, _ := c.Get(key + "/name")
		pubkey, err := PubkeyFromEtcd(c, key+"/pubkey")
		if err != nil {
			return nil, errors.Wrap(err, "Could not get public key")
		}

		conf := KeyConf{
			Name:   name,
			Pubkey: pubkey,
		}

		keys[pubkey.String()] = conf
	}

	return keys, nil
}
