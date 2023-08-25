package cmd

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func RunLoad(cli *cli.Context) error {

	rand.Seed(time.Now().UnixNano())

	nodesCsv := cli.String("nodes")
	nodes := strings.Split(nodesCsv, ",")

	requests := int(cli.Uint("requests"))
	concurrency := int(cli.Uint("concurrency"))

	ch := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	go func() {
		for i := 0; i < requests; i++ {
			wg.Add(1)
			ch <- struct{}{}
		}
	}()

	for i := 0; i < requests; i++ {
		<-ch
		tx := makeTx(nodes[rand.Intn(len(nodes))])
		go sendTx(tx, &wg)
	}

	wg.Wait()

	return nil
}

func makeTx(node string) string {
	tx := make([]byte, 40)
	_, _ = rand.Read(tx)
	hexTx := common.BytesToHash(tx).Hex()
	return fmt.Sprintf("http://%s/broadcast_tx_commit?tx=\"%s\"", node, hexTx)
}

func sendTx(url string, wg *sync.WaitGroup) {

	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("URL: %s, Status Code: %d\n", url, resp.StatusCode)
}
