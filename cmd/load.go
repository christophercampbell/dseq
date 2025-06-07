package cmd

import (
	crand "crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

// RunLoad executes a load test by sending concurrent requests to the specified nodes.
func RunLoad(cli *cli.Context) error {
	nodesCsv := cli.String("nodes")
	if nodesCsv == "" {
		return fmt.Errorf("nodes parameter cannot be empty")
	}
	nodes := strings.Split(nodesCsv, ",")
	if len(nodes) == 0 {
		return fmt.Errorf("no valid nodes provided")
	}

	requests := int(cli.Uint("requests"))
	if requests <= 0 {
		return fmt.Errorf("requests must be greater than 0")
	}

	concurrency := int(cli.Uint("concurrency"))
	if concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}

	// Create a channel to control concurrency
	ch := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	wg.Add(requests)

	start := time.Now()

	// Queue all requests
	go func() {
		for i := 0; i < requests; i++ {
			ch <- struct{}{}
		}
	}()

	var (
		totalDuration time.Duration
		durationMutex sync.Mutex
	)

	stop := make(chan struct{})
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				select {
				case <-ch:
					node := nodes[mrand.Intn(len(nodes))]
					tx := makeTx(node)
					duration := sendTx(tx)

					durationMutex.Lock()
					totalDuration += duration
					durationMutex.Unlock()

					wg.Done()
				case <-stop:
					return
				}
			}
		}()
	}

	wg.Wait()
	close(stop)

	averageRequestTimeMs := float64(totalDuration.Milliseconds()) / float64(requests)
	totalTime := time.Since(start)

	fmt.Printf("Load test results:\n")
	fmt.Printf("  Requests: %d\n", requests)
	fmt.Printf("  Concurrency: %d\n", concurrency)
	fmt.Printf("  Average request time: %.2f ms\n", averageRequestTimeMs)
	fmt.Printf("  Total time: %s\n", totalTime)

	return nil
}

// makeTx generates a random transaction and returns the URL to send it to.
func makeTx(node string) string {
	tx := make([]byte, 40)
	if _, err := crand.Read(tx); err != nil {
		// Fallback to math/rand if crypto/rand fails
		mrand.Read(tx)
	}
	hexTx := common.BytesToHash(tx).Hex()
	return fmt.Sprintf("http://%s/broadcast_tx_commit?tx=\"%s\"", node, hexTx)
}

// sendTx sends a transaction to the specified URL and returns the duration of the request.
func sendTx(url string) time.Duration {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending request to %s: %v\n", url, err)
		return 0
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	// Drain response body
	if _, err := io.ReadAll(resp.Body); err != nil {
		fmt.Printf("Warning: error reading response body: %v\n", err)
	}

	fmt.Printf("URL: %s, Status Code: %d\n", url, resp.StatusCode)
	return time.Since(start)
}
