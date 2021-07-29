// 访问云直连参数
package gwconf

import (
	"github.com/rosbit/go-aes"
	"cdc-gateway/rsa-utils"
	mrand "math/rand"
	"time"
	"fmt"
	"crypto"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"encoding/base64"
)

type CDCConfT struct {
	*AppParams
	privateKey *rsa.PrivateKey
	bankPubKey *rsa.PublicKey
}

var (
	confs = map[string]*CDCConfT{}
)

func AddCDCConf(params *AppParams) {
	confs[params.Name] = &CDCConfT{
		AppParams: params,
	}
}

func GetCDCConf(name string) *CDCConfT {
	if cdcConf, ok := confs[name]; ok {
		return cdcConf
	}
	return nil
}

func (cdcConf *CDCConfT) ChooseUID() string {
	l := len(cdcConf.UIDs)
	if l == 1 {
		return cdcConf.UIDs[0]
	}
	return cdcConf.UIDs[mrand.Intn(l)]
}

func (cdcConf *CDCConfT) parsePrivateKey() error {
	if cdcConf.privateKey == nil {
		pk, err := rsautil.ParsePrivateKeyContent(cdcConf.RSAPrivateKey, "PKCS1")
		if err != nil {
			return err
		}
		cdcConf.privateKey = pk
	}
	return nil
}

func (cdcConf *CDCConfT) parseBankPublicKey() error {
	if cdcConf.bankPubKey == nil {
		pk, err := rsautil.ParsePublicKeyContent(cdcConf.BankPublicKey)
		if err != nil {
			return err
		}
		cdcConf.bankPubKey = pk
	}
	return nil
}

func (cdcConf *CDCConfT) MakeSignature(body map[string]interface{}) (string, error) {
	if err := cdcConf.parsePrivateKey(); err != nil {
		return "", err
	}

	// 增加时间戳
	now := time.Now()
	ts := now.Format("20060102150405")
	body["signature"] = map[string]interface{}{
		"sigtim": ts,
		"sigdat": "__signature_sigdat__",
	}
	bToSign, _ := json.Marshal(body)
	fmt.Printf("strToSign: %s\n", bToSign)
	// Sha256withRSA
	bSha := sha256.Sum256(bToSign)
	signature, err := rsa.SignPKCS1v15(rand.Reader, cdcConf.privateKey, crypto.SHA256, bSha[:])
	if err != nil {
		return "", err
	}
	body["signature"] = map[string]interface{}{
		"sigtim": ts,
		"sigdat": base64.StdEncoding.EncodeToString(signature),
	}
	bToSign, _ = json.Marshal(body)
	fmt.Printf("str with SHA256WithRSA: %s\n", bToSign)
	// AES with ECB mode
	res, err := goaes.AesEncryptECB(bToSign, []byte(cdcConf.AesKey))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func (cdcConf *CDCConfT) ParseResponse(cipher string) (decryptedBody []byte, body map[string]interface{}, err error) {
	if err = cdcConf.parseBankPublicKey(); err != nil {
		return
	}

	res, e := base64.StdEncoding.DecodeString(cipher)
	if e != nil {
		err = e
		return
	}
	bToSign, e := goaes.AesDecryptECB(res, []byte(cdcConf.AesKey))
	if e != nil {
		err = e
		return
	}
	fmt.Printf("decrypted body: %s\n", bToSign)
	decryptedBody = bToSign
	if err = json.Unmarshal(bToSign, &body); err != nil {
		return
	}
	signStruct, ok := body["signature"]
	if !ok || signStruct == nil {
		err = fmt.Errorf("item signature not found")
		return
	}
	var sigtim, sigdat string
	switch signStruct.(type) {
	case map[string]interface{}:
		ss := signStruct.(map[string]interface{})
		sigtimI, ok := ss["sigtim"]
		if !ok {
			err = fmt.Errorf("signature/sigtim not found")
			return
		}
		sigtim, ok = sigtimI.(string)
		if !ok {
			err = fmt.Errorf("string type of signature/sigtim expected")
			return
		}

		sigdatI, ok := ss["sigdat"]
		if !ok {
			err = fmt.Errorf("signature/sigdat not found")
			return
		}
		sigdat, ok = sigdatI.(string)
		if !ok {
			err = fmt.Errorf("string type of signature/sigdat expected")
			return
		}
	default:
		err = fmt.Errorf("bad type for item signature")
		return
	}

	signature, e := base64.StdEncoding.DecodeString(sigdat)
	if e != nil {
		err = e
		return
	}

	body["signature"] = map[string]interface{}{
		"sigtim": sigtim,
		"sigdat": "__signature_sigdat__",
	}
	bToSign, _ = json.Marshal(body)
	fmt.Printf("strToSign: %s\n", bToSign)
	// Sha256withRSA
	bSha := sha256.Sum256(bToSign)
	if err = rsa.VerifyPKCS1v15(cdcConf.bankPubKey, crypto.SHA256, bSha[:], signature); err != nil {
		return
	}
	return
}
