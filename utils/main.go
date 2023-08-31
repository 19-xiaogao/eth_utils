package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"strings"
)

// GenerateAddress 生成私钥和地址
func GenerateAddress(numberAddress int) (privateKeyList []string, addresses []string, err error) {

	for i := 0; i < numberAddress; i++ {
		privateKey, _err := crypto.GenerateKey()
		if _err != nil {
			err = _err
		}
		_privateKeyString := hex.EncodeToString(privateKey.D.Bytes())
		privateKeyList = append(privateKeyList, _privateKeyString)
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			err = errors.New("无法将公钥转换为ECDSA公钥")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
		addresses = append(addresses, address)
	}
	return privateKeyList, addresses, err
}

// PrivateToAddress 私钥导出地址
func PrivateToAddress(private string) (address string) {
	// 将私钥解码为字节数组
	privateKeyBytes, err := hex.DecodeString(private)
	if err != nil {
		log.Fatal(err)
	}

	// 将字节数组转换为私钥
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	// 从私钥中获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法将公钥转换为ECDSA公钥")
	}

	// 生成以太坊地址
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address
}

// SavePrivate 保存私钥到本地文件
func SavePrivate(privateList []string, filePath string) error {
	strData := strings.Join(privateList, "\n")

	err := ioutil.WriteFile(filePath, []byte(strData), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadPrivate 读取本地私钥
func ReadPrivate(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	privateList := strings.Split(string(file), "\n")
	return privateList, nil
}
