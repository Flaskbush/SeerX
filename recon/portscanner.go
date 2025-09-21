package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const target = "192.168.88.112"

func scanPort(target string, port int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp4", address, 2*time.Second)
	if err == nil {
		results <- fmt.Sprintf("Port %d ouvert", port)
		conn.Close()
	}
}

func main() {
	var wg sync.WaitGroup
	results := make(chan string)

	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	for port := 1; port <= 1024; port++ {
		wg.Add(1)
		go scanPort(target, port, &wg, results)
	}

	wg.Wait()
	close(results)
}
