package main

import (
	"context"
	"fmt"
	"time"
)

func withTimeout() {
	duration := time.Second * 2

	ctx, _ := context.WithTimeout(context.Background(), duration)

	select {
		case <- time.After(time.Second * 1):
			fmt.Println("After 1s")
		case <- ctx.Done():
			fmt.Println(ctx.Err())
	}
}

func withCancellation() {
	exit := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Done sleeping 5s. Goodbye")
		cancel()
	}()

	go func() {
		n := 1

		for {
			select {
				case <-ctx.Done():
					fmt.Println("Exiting")
					close(exit)
					return
				default:
					time.Sleep(time.Second)
					fmt.Println("Counting: ", n)
					n += 1
			}
		}
	}()

	<-exit
	fmt.Println("Bye bye")
}

func withValue() {
	ctx := context.WithValue(context.Background(), "jwt", "bearer 12345")

	jwt := ctx.Value("jwt")
	fmt.Println(jwt)
}

func main() {
	// withTimeout()
	// withCancellation()
	// withValue()
}