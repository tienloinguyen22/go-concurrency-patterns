package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
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

func titleize(input <-chan []string) <-chan []string {
	output := make(chan []string)

	go func() {
		for value := range input {
			var titleizedValue []string 
			for _, item := range value {
				item := item
				titleizedValue = append(titleizedValue, strings.Title(item))
			}
			output <- titleizedValue
		}

		close(output)
	}()

	return output
}

func sanitize(input <-chan []string) <-chan []string {
	output := make(chan []string)

	go func() {
		for value := range input {
			removed := false

			for _, item := range value {
				if len(item) > 3 {
					removed = true
				}
			}

			if !removed {
				output <- value
			}
		}

		close(output)
	}()

	return output
}

func main() {
	readChan, err := read("file.csv")
	if err != nil {
		fmt.Printf("Error reading file %v\n", err)
		return
	}

	for value := range sanitize(titleize(readChan)) {
		fmt.Println(value)
	}
}