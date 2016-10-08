package curveface

import "github.com/pkg/errors"

// GetGface quick gface request
func GetGface(endpoint, priv, pub, srvpub string) ([12]byte, error) {
	noop := [12]byte{}
	ssocket, err := Client{}.Dial(endpoint, priv, pub, srvpub)
	if err == nil {
		_, err = ssocket.Write([]byte("gface"))
		if err != nil {
			return noop, errors.Wrap(err, "Failed to write to wrapped socket")
		}

		var packet [12]byte

		_, err = ssocket.Read(packet[:])
		if err != nil {
			return noop, errors.Wrap(err, "Failed to read from wrapped socket")
		}

		err = ssocket.Close()
		if err != nil {
			return noop, errors.Wrap(err, "Failed to close socket")
		}

		return packet, nil
	}
	return noop, err
}
