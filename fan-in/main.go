package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

func read(filename string) (<-chan []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not open file %v", err)
	}

	fileReader := csv.NewReader(file)
	ch := make(chan []string)

	go func() {
		for {
			row, err := fileReader.Read()
			if err == io.EOF {
				close(ch)
				return
			}
			ch <- row
		}
	}()

	return ch, nil
}

func mergeChannelsUsingWaitGroup(chans ...<-chan []string) <-chan []string {
	var wg sync.WaitGroup
	outputChan := make(chan []string)

	send := func(ch <-chan []string) {
		for value := range ch {
			outputChan <- value
		}

		wg.Done()
	}

	wg.Add(len(chans))
	for _, ch := range chans {
		go send(ch)
	}


	go func() {
		wg.Wait()
		close(outputChan)
	}()

	return outputChan
}

func main() {
	chan1, err := read("file1.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file1 %v", err))
	}

	chan2, err := read("file2.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file2 %v", err))
	}

	mergedChan := mergeChannelsUsingWaitGroup(chan1, chan2)

	for value := range mergedChan {
		fmt.Println(value)
	}

	fmt.Println("All completed, exiting!")
}