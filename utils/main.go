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

// GenerateAddress Generate the private key and address
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
			err = errors.New("cannot convert a public key to an ECDSA public key")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
		addresses = append(addresses, address)
	}
	return privateKeyList, addresses, err
}

// PrivateToAddress Private key export address
func PrivateToAddress(private string) (address string) {

	privateKeyBytes, err := hex.DecodeString(private)
	if err != nil {
		log.Fatal(err)
	}

	// Converts the byte array into a private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Get the public key from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法将公钥转换为ECDSA公钥")
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address
}

// SavePrivate Save the private key to a local file
func SavePrivate(privateList []string, filePath string) error {
	strData := strings.Join(privateList, "\n")

	err := ioutil.WriteFile(filePath, []byte(strData), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadLocalPrivate Read the local private key
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
