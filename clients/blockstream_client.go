package clients

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"transaction-ancestory/models"
)

var GetHashPath = "block-height/%s"
var GetTransactionsPath = "block/%s/txs/%d"

type BlockStreamClient struct {
	host       string
	httpClient http.Client
}

type BlockStreamClientInterface interface {
	GetHash(block string) (string, error)
	GetTransactions(hash string, startIndex int) (*models.HttpResponse, error)
}

func NewBlockStreamClient() *BlockStreamClient {
	return &BlockStreamClient{
		host:       "https://blockstream.info/api/",
		httpClient: http.Client{},
	}
}

func (bsc *BlockStreamClient) GetHash(block string) (string, error) {
	path := fmt.Sprintf(GetHashPath, block)
	response, err := bsc.httpClient.Get(fmt.Sprintf("%s%s", bsc.host, path))
	if err != nil {
		fmt.Println("Failed to fetch hash", err.Error())
		return "", errors.New("Failed to fetch hash")
	}
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}

func (bsc *BlockStreamClient) GetTransactions(hash string, startIndex int) (*models.HttpResponse, error) {
	httpResponse := models.HttpResponse{}
	path := fmt.Sprintf(GetTransactionsPath, hash, startIndex)
	response, err := bsc.httpClient.Get(fmt.Sprintf("%s%s", bsc.host, path))
	if err != nil {
		fmt.Println("Failed to fetch hash", err.Error())
		return &httpResponse, errors.New("Failed to fetch hash")
	}
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &httpResponse, err
	}

	httpResponse.Body = respBytes
	httpResponse.Status = response.StatusCode

	return &httpResponse, nil
}
