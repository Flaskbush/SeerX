package recon

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct{}

func NewPortScanner() *PortScanner {
	return &PortScanner{}
}

func (ps *PortScanner) Scan(target, portRange string) ([]string, error) {
	var wg sync.WaitGroup
	results := make(chan string)
	var openPorts []string

	go func() {
		for port := range results {
			openPorts = append(openPorts, port)
		}
	}()

	ports := parsePortRange(portRange)
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err == nil {
				results <- fmt.Sprintf("Port %d open", p)
				conn.Close()
			}
		}(port)
	}

	wg.Wait()
	close(results)
	return openPorts, nil

}

func parsePortRange(portRange string) []int {
	var ports []int
	parts := strings.Split(portRange, "-")
	start, _ := strconv.Atoi(parts[0])
	end := start
	if len(parts) > 1 {
		end, _ = strconv.Atoi(parts[1])
	}
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return ports
}
