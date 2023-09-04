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

// Distribute 分发: 一个地址向多个地址发送eth
func (server *Server) Distribute(privateKey string, address []string, amount float64) error {
	// 私钥推算地址
	fromAddress := utils.PrivateToAddress(privateKey)
	// 计算该私钥的余额是否可以分发到怎么多地址
	balance, err := server.client.BalanceAt(context.Background(), common.HexToAddress(fromAddress), nil)
	if err != nil {
		return err
	}
	totalAmount := utils.ToWei(amount*float64(len(address)), 18)

	if balance.Cmp(totalAmount) < 0 {
		return errors.New("该地址的余额不够分发这么多地址")
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

	// 批量发送
}
