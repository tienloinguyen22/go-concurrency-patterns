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

func breakup(workerId string, ch <-chan []string) {
	for value := range ch {
		fmt.Println(workerId, value)
	}
}

func main() {
	ch, err := read("file.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file1 %v", err))
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		breakup("1", ch)
		wg.Done()
	}()
	go func() {
		breakup("2", ch)
		wg.Done()
	}()
	go func() {
		breakup("3", ch)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("All completed, exiting!")
}