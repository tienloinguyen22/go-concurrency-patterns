package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
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

func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	filenames := []string{"file1.csv", "file2.csv", "file3.csv"}
	for _, filename := range filenames {
		filename := filename
		eg.Go(func() error {
			ch, err := read(filename)
			if err != nil {
				return fmt.Errorf("Error reading %v\n", err)
			}
			if filename == "file1.csv" {
				time.Sleep(100 * time.Microsecond)
				return fmt.Errorf("Random error after 50ms\n") 
			}

			for {
				select {
					case <-ctx.Done():
						fmt.Printf("Context completed from %v. Error: %v\n", filename, ctx.Err())
						return ctx.Err() 
					case value, open := <-ch:
						if !open {
							return nil
						}
						fmt.Println(value)
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error reading files: %v\n", err)
	}
}