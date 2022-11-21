package rsa_crypt

import (
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/convert"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"math/big"
)

// ukey 导出base64字串转换成pem格式
func Translate(source string) string {
	keyByte, _ := base64.StdEncoding.DecodeString(source)
	bits := new(big.Int)
	bits.SetBytes(utils.ReverseBytes(keyByte[:4]))
	e := new(big.Int)
	e.SetBytes(utils.ReverseBytes(keyByte[4:8]))
	n := new(big.Int)
	readLen := bits.Int64()/8 + 8
	n.SetBytes(keyByte[8:readLen])
	pubKey := rsa.PublicKey{E: convert.ToInt(e.Int64()), N: n}
	publicStream := x509.MarshalPKCS1PublicKey(&pubKey)
	return base64.StdEncoding.EncodeToString(publicStream)
}
