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

func RunLoad(cli *cli.Context) error {
	nodesCsv := cli.String("nodes")
	nodes := strings.Split(nodesCsv, ",")

	requests := int(cli.Uint("requests"))
	concurrency := int(cli.Uint("concurrency"))

	ch := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	wg.Add(requests)

	start := time.Now()

	go func() {
		for i := 0; i < requests; i++ {
			ch <- struct{}{}
		}
	}()

	totalDuration := 0 * time.Second

	stop := make(chan struct{})
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				select {
				case <-ch:
					tx := makeTx(nodes[mrand.Intn(len(nodes))])
					t := sendTx(tx)
					wg.Done()
					totalDuration += t
				case <-stop:
					return
				}

			}
		}()
	}

	wg.Wait()
	stop <- struct{}{}

	averageRequestTimeMs := float64(totalDuration.Milliseconds()) / float64(requests)

	fmt.Printf("requests %d, concurrency %d, avg: %v ms, total %s\n", requests, concurrency, averageRequestTimeMs, time.Now().Sub(start))
	return nil
}

func makeTx(node string) string {
	tx := make([]byte, 40)
	_, _ = crand.Read(tx)
	hexTx := common.BytesToHash(tx).Hex()
	return fmt.Sprintf("http://%s/broadcast_tx_commit?tx=\"%s\"", node, hexTx)
}

func sendTx(url string) time.Duration {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	defer func() { _ = resp.Body.Close() }()

	// drain response as a real client would do, but just warn on error
	if _, warn := io.ReadAll(resp.Body); warn != nil {
		fmt.Println("Warn:", err)
	}

	fmt.Printf("URL: %s, Status Code: %d\n", url, resp.StatusCode)
	return time.Now().Sub(start)
}
