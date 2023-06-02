package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"transaction-ancestory/clients"
	"transaction-ancestory/models"
)

type BlockStreamService struct {
	blockStreamClient clients.BlockStreamClientInterface
}

type BlockStreamServiceInterface interface {
	GetAncestorCount(block string) ([]models.AncestorCount, error)
}

func NewBlockStreamService() *BlockStreamService {
	return &BlockStreamService{
		blockStreamClient: clients.NewBlockStreamClient(),
	}
}

func (bss *BlockStreamService) GetAncestorCount(block string) ([]models.AncestorCount, error) {

	hash, err := bss.blockStreamClient.GetHash(block)
	if err != nil {
		fmt.Println("Faile to fetch hash from block", block)
		return nil, errors.New("Faile to fetch hash from block")
	}

	transactions, err := bss.fetchAllTransactions(hash)
	if err != nil {
		return nil, errors.New("Failed to fetch the transactions")
	}

	ancestorCountMap := bss.calculateCount(transactions)
	ancestorCounts := bss.findTopN(ancestorCountMap, 10)

	return ancestorCounts, nil
}

func (bss *BlockStreamService) fetchAllTransactions(hash string) ([]*models.Transaction, error) {
	transactions := []*models.Transaction{}
	startIndex := 0
	for {
		httpResponse, err := bss.blockStreamClient.GetTransactions(hash, startIndex)
		if err != nil {
			fmt.Println("Failed to unmarshal response to json", err.Error())
			return transactions, errors.New("Failed to unmarshal response to json")
		}

		if httpResponse.Status == 404 {
			fmt.Println("Response status 404")
			break
		}

		if httpResponse.Status != 200 {
			fmt.Println("Response status is not 200", httpResponse.Status)
			return transactions, errors.New("Error while fetching the transactions")
		}

		newTransactions := []*models.Transaction{}
		err = json.Unmarshal(httpResponse.Body, &newTransactions)
		transactions = append(transactions, newTransactions...)
		startIndex += 25
		fmt.Println("Fetching transaction index", startIndex)
	}

	return transactions, nil
}

func (bss *BlockStreamService) findTopN(ancestorCountMap map[string]int, N int) []models.AncestorCount {
	ancestorsCounts := []models.AncestorCount{}

	for key, value := range ancestorCountMap {
		ancestorsCounts = append(ancestorsCounts, models.AncestorCount{
			Txid:  key,
			Count: value,
		})
	}

	sort.SliceStable(ancestorsCounts, func(i, j int) bool {
		return ancestorsCounts[i].Count > ancestorsCounts[j].Count
	})

	return ancestorsCounts[:N]
}

func (bss *BlockStreamService) calculateCount(transactions []*models.Transaction) map[string]int {
	ancestorsCount := map[string]int{}
	parentsCount := map[string]int{}
	parentsMap := map[string][]string{}
	txidMap := map[string]bool{}

	for _, tx := range transactions {
		txidMap[tx.Txid] = true
	}

	for i, tx := range transactions {
		for _, inputTx := range tx.InputTxn {
			if txidMap[inputTx.Txid] && transactions[i].Txid != inputTx.Txid {
				parentsCount[transactions[i].Txid] += 1
				parentsMap[transactions[i].Txid] = append(parentsMap[transactions[i].Txid], inputTx.Txid)
			}
		}
	}

	for txid := range parentsMap {
		queue := []string{txid}
		for len(queue) > 0 {
			ancestorsCount[txid] += parentsCount[queue[0]]
			for ind := range parentsMap[queue[0]] {
				queue = append(queue, parentsMap[queue[0]][ind])
			}

			queue = queue[1:]
		}
	}

	return ancestorsCount
}

func findTransactionByID(transactions []*models.Transaction, txid string) *models.Transaction {
	for _, tx := range transactions {
		if tx.Txid == txid {
			return tx
		}
	}
	return &models.Transaction{}
}
