package main

import (
	"fmt"
	"transaction-ancestory/services"
)

func main() {
	fmt.Println("Enter the block number")
	block := ""
	fmt.Scanf("%s", &block)

	blockStreamService := services.NewBlockStreamService()
	counts, err := blockStreamService.GetAncestorCount(block)
	if err != nil {
		fmt.Println("Error while fetching ancetors count", err.Error())
		return
	}

	for _, count := range counts {
		fmt.Printf("Txn id: %s, Count: %d\n", count.Txid, count.Count)
	}
}
