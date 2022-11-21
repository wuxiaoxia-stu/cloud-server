package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	FrameSize = 2048

	FPFileVersion1 uint16 = 0xfe
	FPFileVersion2 uint16 = 0xef
	FPFileVersion3 uint8  = 1

	FDInfo string = "蒹葭苍苍，白露为霜。所谓伊人，在水一方。"
)

//只能16、24、32共三种情况
func pwdPadding(pwd []byte) []byte {
	maxSize := 32
	pwdSize := len(pwd)
	if pwdSize == 32 {
		return pwd
	} else if pwdSize > 32 {
		return pwd[0:32]
	} else {
		padding := maxSize - pwdSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(pwd, padtext...)
	}
	// if pwdSize%minSize == 0 {
	// 	return pwd
	// } else if pwdSize <= minSize {
	// 	padding := minSize - pwdSize
	// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// 	return append(pwd, padtext...)
	// } else {
	// 	padding := minSize - pwdSize%minSize
	// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// 	return append(pwd, padtext...)
	// }
}

func EncryptFile(src, dest, pwd string) error {

	correctPwd := pwdPadding([]byte(pwd))

	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	//写文件头
	var encryptedInfo = AesEncryptCBC([]byte(FDInfo), correctPwd)
	var fdHeader []byte = make([]byte, 6)
	binary.BigEndian.PutUint16(fdHeader[0:], FPFileVersion1)
	binary.BigEndian.PutUint16(fdHeader[2:], FPFileVersion2)
	binary.BigEndian.PutUint16(fdHeader[4:], uint16(len(encryptedInfo)))
	destfile.Write(fdHeader)
	destfile.Write(encryptedInfo)

	//加密处理
	frameBuf := make([]byte, FrameSize) //一次读取多少个字节
	sizeBuf := make([]byte, 2)
	for {
		n, err := srcfile.Read(frameBuf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
		if n <= 0 {
			break
		}

		encryptedFrame := AesEncryptCBC(frameBuf[:n], correctPwd)
		binary.BigEndian.PutUint16(sizeBuf[:], uint16(len(encryptedFrame)))
		_, err = destfile.Write(sizeBuf)
		if err != nil {
			return err
		}
		_, err = destfile.Write(encryptedFrame)
		if err != nil {
			return err
		}
	}

	return nil
}

func DecryptFile(src, dest, pwd string) error {
	correctPwd := pwdPadding([]byte(pwd))

	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	//读文件头
	fdHeader := make([]byte, 6)
	srcfile.Read(fdHeader)

	v1 := binary.BigEndian.Uint16(fdHeader[0:])
	v2 := binary.BigEndian.Uint16(fdHeader[2:])
	encryptedInfoSize := binary.BigEndian.Uint16(fdHeader[4:])

	if v1 != FPFileVersion1 && v2 != FPFileVersion2 {
		return errors.New("Not a FP file")
	}
	encryptedInfo := make([]byte, encryptedInfoSize)
	srcfile.Read(encryptedInfo)

	descryptedFDInfo := AesEncryptCBC([]byte(FDInfo), correctPwd)

	if strings.Compare(string(descryptedFDInfo), string(encryptedInfo)) != 0 {
		return errors.New("password wrong")
	}

	sizeBuf := make([]byte, 2)
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()
	for {
		n, err := srcfile.Read(sizeBuf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
		if n <= 0 {
			break
		}
		encryptedFrameSize := binary.BigEndian.Uint16(sizeBuf)
		buf := make([]byte, encryptedFrameSize)
		n, err = srcfile.Read(buf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}

		if n <= 0 {
			break
		}

		d := AesDecryptCBC(buf[:n], correctPwd)
		destfile.Write(d)
	}

	return nil
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	newCipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, newCipher.BlockSize(); bs <= len(origData); bs, be = bs+newCipher.BlockSize(), be+newCipher.BlockSize() {
		newCipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	newCipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))

	for bs, be := 0, newCipher.BlockSize(); bs < len(encrypted); bs, be = bs+newCipher.BlockSize(), be+newCipher.BlockSize() {
		newCipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
