package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortResult struct {
	Port    int
	Open    bool
	Service string
	Error   error
}

func ScanPort(host string, port int, timeout time.Duration) PortResult {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return PortResult{Port: port, Open: false, Error: err}
	}
	defer conn.Close()
	return PortResult{Port: port, Open: true, Error: nil}
}

func ScanPorts(host string, ports []int, timeout time.Duration, concurrency int) []PortResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]PortResult, 0, len(ports))
	portsChan := make(chan int, len(ports))

	//feed ports to channel
	for _, port := range ports {
		portsChan <- port
	}
	close(portsChan) //tells goroutines no more ports are incoming

	//start worker pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsChan {
				result := ScanPort(host, port, timeout)

				if result.Open {
					fmt.Printf("[+] Port %d open\n", result.Port)
				} else if result.Error != nil {
					fmt.Printf("[-] Port %d error: %v\n", result.Port, result.Error)
				} else {
					fmt.Printf("[-] Port %d open\n", result.Port)
				}
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}()
	}
	wg.Add(-1)
	return results
}
