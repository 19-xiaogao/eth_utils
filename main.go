package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"mass_address/serve"
	"mass_address/utils"
	"os"
	"strconv"
)

func init() {
	fmt.Println("use generateAddress 10 create address")
	fmt.Println("use savePrivate 10 ./priv.text save created private key")
	fmt.Println("use readPrivate ./priv.text read create private key")
	fmt.Println("use privateToAddress adfb export private to address")
	fmt.Println("use distribute [privateKey] ./address.text 10")
	fmt.Println("use collection ./pri.text oxasdfaf")

}

// PRC
const RPC = "https://eth.getblock.io/6a3095cc-5bf1-4977-b7a9-6d5e5de64ca3/goerli/"

func main() {
	client, err := ethclient.Dial(RPC)
	if err != nil {
		log.Fatal(err)
	}
	serverInterface := serve.NewServer(client)
	//err = serverInterface.SendEthTx(privateKey, recipientAddress, amount)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// Uncomment this block to pass the first stage!
	//
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	switch command := os.Args[1]; command {

	case "distribute":
		prvKey := os.Args[2]
		filePath := os.Args[3]
		amount, _ := strconv.ParseFloat(os.Args[4], 64)
		addressList, err := utils.ReadLocalPrivate(filePath)
		if err != nil {
			log.Fatal(err)
		}
		err = serverInterface.Distribute(prvKey, addressList, amount)
		if err != nil {
			log.Fatal(err)
		}
	case "collection":
		filePath := os.Args[2]
		collectionAddress := os.Args[3]
		provList, err := utils.ReadLocalPrivate(filePath)
		if err != nil {
			log.Fatal(err)
		}
		err = serverInterface.Collection(provList, collectionAddress)
		if err != nil {
			log.Fatal(err)
		}

	case "generateAddress":
		num, err := strconv.Atoi(os.Args[2])
		privateKeyList, addresses, err := utils.GenerateAddress(num)
		if err != nil {
			return
		}
		for _, value := range privateKeyList {
			fmt.Println(value)
		}
		fmt.Println()
		for _, value := range addresses {
			fmt.Println(value)
		}
		fmt.Println()
	case "savePrivate":
		filepath := os.Args[3]
		num, err := strconv.Atoi(os.Args[2])
		privateKeyList, _, err := utils.GenerateAddress(num)
		err = utils.SavePrivate(privateKeyList, filepath)
		if err != nil {
			return
		}
		fmt.Println("save success")

	case "readPrivate":
		filePath := os.Args[2]
		privateKeyList, err := utils.ReadLocalPrivate(filePath)
		if err != nil {
			return
		}
		fmt.Println(privateKeyList)
	case "privateToAddress":
		private := os.Args[2]
		fmt.Println("address: ", utils.PrivateToAddress(private))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
