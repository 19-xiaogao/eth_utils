package main

import (
	"fmt"
	"mass_address/utils"
	"os"
	"strconv"
)

func init() {
	fmt.Println("通过 GenerateAddress 10 创建地址")
	fmt.Println("通过 savePrivate 10 ./priv.text 保存创建的私钥")
	fmt.Println("通过 readPrivate ./priv.text 读取创建的私钥")
	fmt.Println("通过 privateToAddress adfb 推倒出私钥")
}
func main() {

	// Uncomment this block to pass the first stage!
	//
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	switch command := os.Args[1]; command {

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
		privateKeyList, err := utils.ReadPrivate(filePath)
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
