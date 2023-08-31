package main

import (
	"fmt"
	"mass_address/utils"
)

func main() {
	//privateKeyList, addresses, err := utils.GenerateAddress(100)
	//if err != nil {
	//	return
	//}

	//err = utils.SavePrivate(privateKeyList, "private.txt")
	//if err != nil {
	//	return
	//}

	privateKeyList, err := utils.ReadPrivate("private.txt")
	if err != nil {
		return
	}
	fmt.Println(privateKeyList)
	//for _, value := range privateKeyList {
	//	fmt.Println(value)
	//}
	//fmt.Println()
	//for _, value := range addresses {
	//	fmt.Println(value)
	//}
	//fmt.Println()

	//fmt.Println(utils.PrivateToAddress(privateKeyList[0]))

}
