package main

import (
	"gateway/grpc"
	"gateway/proxy"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		proxy.Start()
		wg.Done()
	}()

	go func() {
		grpc.Start()
		wg.Done()
	}()

	wg.Wait()

}
