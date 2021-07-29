package rsautil

import (
	"crypto/rsa"
	"io/ioutil"
)

func ParsePrivateKey(priKeyFile string, typ string) (pri *rsa.PrivateKey, err error) {
	b, e := ioutil.ReadFile(priKeyFile)
	if e != nil {
		err = e
		return
	}

	return parsePrivateKey(b, typ)
}

func ParsePublicKey(pubKeyFile string) (pub *rsa.PublicKey, err error) {
	b, e := ioutil.ReadFile(pubKeyFile)
	if e != nil {
		err = e
		return
	}

	return parsePublicKey(b)
}
