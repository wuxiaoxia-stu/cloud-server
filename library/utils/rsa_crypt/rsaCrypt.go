package rsa_crypt

import (
	"aiyun_cloud_srv/library/utils/convert"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
)

type RsaKeyPair struct {
	PubKey string
	PriKey string
}

func DumpPrivateKeyBase64(privatekey *rsa.PrivateKey) (string, error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

func DumpPublicKeyBase64(publickey *rsa.PublicKey) (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		g.Log().Error(err)
		return "", err
	}

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

// Load private key from base64
func LoadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		g.Log().Error(err)
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		g.Log().Error(err)
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}

func LoadPublicKeyBase64(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		g.Log().Error(err)
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	publickey, err := x509.ParsePKCS1PublicKey(keybytes)
	if err != nil {
		g.Log().Error(err)
		return nil, err
	}
	return publickey, nil
}

func RSAGenKey(bits int, rsaKeyPair *RsaKeyPair) (err error) {
	/*
		生成私钥
	*/
	//1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		g.Log().Error(err)
		return
	}
	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	rsaKeyPair.PriKey = base64.StdEncoding.EncodeToString(privateStream)
	publicKey := privateKey.PublicKey
	publicStream := x509.MarshalPKCS1PublicKey(&publicKey)
	rsaKeyPair.PubKey = base64.StdEncoding.EncodeToString(publicStream)
	return nil
}

func Crypt(pubKey *rsa.PublicKey, message string) string {
	rng := rand.Reader
	hashed := []byte(message)
	cipher, err := rsa.EncryptPKCS1v15(rng, pubKey, hashed[:])
	if err != nil {
		g.Log().Error(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(cipher)
}

func DeCrypt(priKey *rsa.PrivateKey, cipher string) string {
	rng := rand.Reader
	cipherByte, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		g.Log().Error(err)
		return ""
	}
	plainTextByte, err := rsa.DecryptPKCS1v15(rng, priKey, cipherByte)
	if err != nil {
		g.Log().Error(err)
		return ""
	}
	return convert.ToString(plainTextByte)
}

func Sign(priKey *rsa.PrivateKey, message string, isURLEncode ...bool) string {
	rng := rand.Reader
	hashed := []byte(message)
	signature, err := rsa.SignPKCS1v15(rng, priKey, 0, hashed[:])
	if err != nil {
		g.Log().Error(err)
		return ""
	}
	if len(isURLEncode) > 0 && isURLEncode[0] {
		return base64.URLEncoding.EncodeToString(signature)
	}
	return base64.StdEncoding.EncodeToString(signature)
}

func Verify(pubKey *rsa.PublicKey, message, signature string, isURLEncode ...bool) bool {
	var sign []byte
	var err error
	if len(isURLEncode) > 0 && isURLEncode[0] {
		sign, err = base64.URLEncoding.DecodeString(signature)
	} else {
		sign, err = base64.StdEncoding.DecodeString(signature)
	}
	if err != nil {
		g.Log().Error(err)
		return false
	}

	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed. This requires
	// that the hash function be collision resistant. SHA-256 is the
	// least-strong hash function that should be used for this at the time
	// of writing (2016).
	hashed := []byte(message)
	err = rsa.VerifyPKCS1v15(pubKey, 0, hashed[:], sign)
	if err != nil {
		g.Log().Error(err)
		return false
	}
	return true
}
