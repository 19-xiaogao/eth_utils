package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"math/big"
	"regexp"
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

// ReadLocalPrivate 读取本地私钥
func ReadLocalPrivate(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	privateList := strings.Split(string(file), "\n")
	return privateList, nil
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}
