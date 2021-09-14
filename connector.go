package bitcoindata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIConnector struct {
	fetchedBlockHeight int
}

func NewAPIConnector(initialBlockHeight int) *APIConnector {

	if initialBlockHeight < 0 {
		panic(fmt.Errorf("initialBlock should be bigger than or equal to 0"))
	}

	apiConnector := &APIConnector{fetchedBlockHeight: initialBlockHeight}
	return apiConnector
}

func fetch(blockHeight int) ([]Block, error) {

	response, err := http.Get(fmt.Sprintf("https://blockchain.info/block-height/%d", blockHeight))

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Blocks, nil
}

func (a *APIConnector) FetchBlock(blockHeight int) ([]Block, error) {
	return fetch(blockHeight)
}

func (a *APIConnector) FetchNextHeightBlocks() ([]Block, error) {

	blocks, err := fetch(a.fetchedBlockHeight)
	a.fetchedBlockHeight++
	return blocks, err
}
