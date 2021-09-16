package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/korkmazkadir/bitcoindata"
)

func main() {

	initialRount := 600000
	fetchCount := 10000
	concurrency := 100

	//blockChan := make(chan bitcoindata.Block, fetchCount)

	workChanCappacity := fetchCount / concurrency

	var workChans []chan int
	var blockChans []chan bitcoindata.Block

	// start workers
	var wgWorker sync.WaitGroup
	for i := 0; i < concurrency; i++ {

		workChan := make(chan int, workChanCappacity)
		workChans = append(workChans, workChan)

		blockChan := make(chan bitcoindata.Block, 2)
		blockChans = append(blockChans, blockChan)

		wgWorker.Add(1)
		go func() {
			defer wgWorker.Done()
			getData(workChan, blockChan)
		}()
	}

	// give works
	for i := initialRount; i < initialRount+fetchCount; i++ {
		ch := workChans[i%len(workChans)]
		ch <- i
	}

	// start workers
	var wgPersist sync.WaitGroup
	wgPersist.Add(1)
	go func() {
		defer wgPersist.Done()
		persist(blockChans)
	}()

	// wait for workers to finish
	wgWorker.Wait()

	// wait for persist work to finish
	wgPersist.Wait()

}

func persist(blockChans []chan bitcoindata.Block) {
	csvFile, err := os.OpenFile("block_header.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer csvFile.Close()

	if _, err := csvFile.WriteString(bitcoindata.BlockCSVHeader()); err != nil {
		panic(err)
	}

	index := 0
	for {

		block, more := <-blockChans[index%len(blockChans)]

		if !more {
			return
		}

		fmt.Printf("[%d] %s\n", block.BlockIndex, block.Hash)
		if _, err := csvFile.WriteString(block.CSVString()); err != nil {
			panic(err)
		}

		index++
	}

}

func getData(workChan chan int, blockChan chan bitcoindata.Block) {

	connector := bitcoindata.NewAPIConnector(0)

	for {
		blockHeight, more := <-workChan

		if !more {
			return
		}

		var blocks []bitcoindata.Block
		var err error

		for {
			blocks, err = connector.FetchBlock(blockHeight)

			if err == nil {
				break
			}

			time.Sleep(2 * time.Second)
		}

		if len(blocks) == 1 {

			blockChan <- blocks[0]
		} else if len(blocks) > 1 {

			fmt.Printf("[%d] number of blocks is %d\n", blockHeight, len(blocks))

			for _, block := range blocks {
				blockChan <- block
			}

		}

	}

}
