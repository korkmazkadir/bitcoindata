package main

import (
	"fmt"
	"sync"

	"github.com/korkmazkadir/bitcoindata"
)

func main() {

	initialRount := 200000
	fetchCount := 10000
	concurrency := 10

	//blockChan := make(chan bitcoindata.Block, fetchCount)

	workChanCappacity := fetchCount / concurrency

	var workChans []chan int
	var blockChans []chan bitcoindata.Block

	// start workers
	var wgWorker sync.WaitGroup
	for i := 0; i < concurrency; i++ {

		workChan := make(chan int, workChanCappacity)
		workChans = append(workChans, workChan)

		blockChan := make(chan bitcoindata.Block, 3)
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

	index := 0
	for {

		block, more := <-blockChans[index%len(blockChans)]

		if !more {
			return
		}

		fmt.Printf("[%d] %s\n", block.BlockIndex, block.Hash)

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

		blocks, err := connector.FetchBlock(blockHeight)

		if err != nil {
			panic(err)
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
