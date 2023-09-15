package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

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
