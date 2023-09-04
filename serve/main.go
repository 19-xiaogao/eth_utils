package serve

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"mass_address/utils"
	"math/big"
)

type Server struct {
	client *ethclient.Client
}

func NewServer(client *ethclient.Client) *Server {
	return &Server{client: client}
}

// SendEthTx 使用私钥发送交易
func (server *Server) SendEthTx(privateKey string, recipientAddress string, amount interface{}) error {

	_privateKey, _ := crypto.HexToECDSA(privateKey)

	fromAddress := utils.PrivateToAddress(privateKey)

	_amount := utils.ToWei(amount, 18)

	balance, err := server.client.BalanceAt(context.Background(), common.HexToAddress(fromAddress), nil)
	if err != nil {
		return err
	}
	if _amount.Cmp(balance) > 0 {
		log.Fatal("余额不足")
	}
	// 获取地址的nonce
	nonce, err := server.client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))

	if err != nil {
		return err
	}

	gasLimit := uint64(21000)

	gasPrice, err := server.client.SuggestGasPrice(context.Background())

	if err != nil {
		return err

	}

	toAddress := common.HexToAddress(recipientAddress)
	tx := types.NewTransaction(nonce, toAddress, _amount, gasLimit, gasPrice, nil)
	chainID, err := server.client.NetworkID(context.Background())
	if err != nil {
		return err

	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), _privateKey)
	if err != nil {
		return err

	}
	// 发送交易
	err = server.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err

	}
	fmt.Printf("交易已发送，交易哈希：%s\n", signedTx.Hash().Hex())

	return nil
}

// Distribute  One address sends eth to multiple addresses
func (server *Server) Distribute(privateKey string, address []string, amount float64) error {

	fromAddress := utils.PrivateToAddress(privateKey)

	balance, err := server.client.BalanceAt(context.Background(), common.HexToAddress(fromAddress), nil)
	if err != nil {
		return err
	}
	totalAmount := utils.ToWei(amount*float64(len(address)-1), 18)

	if balance.Cmp(totalAmount) < 0 {
		return errors.New("the balance of this address is not enough to distribute so many addresses")
	}
	for i := 0; i < len(address)-1; i++ {
		if !utils.IsValidAddress(address[i]) {
			fmt.Println("")
			continue
		}
		//go func() {
		//	err := server.SendEthTx(privateKey, address[i], amount)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//}()
		err := server.SendEthTx(privateKey, address[i], amount)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

// Collection Send the eth of all private keys to the destination address
func (server *Server) Collection(privateList []string, collectionAddress string) error {
	gasLimit := uint64(21000)

	gasPrice, _ := server.client.SuggestGasPrice(context.Background())
	gasVal := big.NewInt(0).Mul(gasPrice, big.NewInt(0).SetUint64(gasLimit))
	for i := 0; i < len(privateList)-1; i++ {

		fromAddress := utils.PrivateToAddress(privateList[i])
		balance, err := server.client.BalanceAt(context.Background(), common.HexToAddress(fromAddress), nil)
		if err != nil {
			return err
		}
		allBalance := big.NewInt(0).Sub(balance, gasVal)
		err = server.SendEthTx(privateList[i], collectionAddress, utils.ToDecimal(allBalance, 18))
		if err != nil {
			return err
		}
	}
	return nil
}
