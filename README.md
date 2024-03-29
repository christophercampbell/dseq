# desq

A decentralized byte sequencer built with cometbft & zkevm-data-stream that produces a uniform sequence of bytes from multiple uncoordinated producers. 

![BFT sequencer](dseq.png)

### getting started

Install cometbft, it has tools for generating node configs.

```shell
git clone --depth 1 --branch v0.38.2  https://github.com/cometbft/cometbft.git
cd cometbft
make install
```

### build the binary

```shell
make build
```

Cometbft generates a test cluster of nodes that are configured as peers. The docker compose file defines 4 nodes that use these configurations.

To run the local testnet cluster

```shell
make start
```

Send a transaction

```
curl -s 'localhost:26657/broadcast_tx_commit?tx="0xDEADBEEF"'
```

Send in some load randomly across nodes

```shell
make load-10                                                                                                                                                                                       
URL: http://localhost:26657/broadcast_tx_commit?tx="0x8b763a1b1d49d4955c8486216325253fec738dd7a9e28bf921119c160f070244", Status Code: 200
URL: http://localhost:26657/broadcast_tx_commit?tx="0xb14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b", Status Code: 200
URL: http://localhost:26660/broadcast_tx_commit?tx="0xbdf8883a0ad8be9c3978b04883e56a156a8de563afa467d49dec6a40e9a1d007", Status Code: 200
URL: http://localhost:26660/broadcast_tx_commit?tx="0x7c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939cb6694d2c422ac", Status Code: 200
URL: http://localhost:26657/broadcast_tx_commit?tx="0xd85794bb358b0c3b525da1786f9fff094279db1944ebd7a19d0f7bbacbe0255a", Status Code: 200
URL: http://localhost:26662/broadcast_tx_commit?tx="0xb0223beea5f4f74391f445d15afd4294040374f6924b98cbf8713f8d962d7c8d", Status Code: 200
URL: http://localhost:26662/broadcast_tx_commit?tx="0x7a4c7215a3b539eb1e5849c6077dbb5722f5717a289a266f97647981998ebea8", Status Code: 200
URL: http://localhost:26662/broadcast_tx_commit?tx="0x9875921e668a5bdf2c7fc4844592d2572bcd0668d2d6c52f5054e2d0836bf84c", Status Code: 200
URL: http://localhost:26664/broadcast_tx_commit?tx="0x7311e4d7defa922daae7786667f7e936cd4f24abf7df866baa56038367ad6145", Status Code: 200
URL: http://localhost:26664/broadcast_tx_commit?tx="0x84045d87f3c67cf22746e995af5a25367951baa2ff6cd471c483f15fb90badb3", Status Code: 200
```

Compare the node's sequence files, they should be the same. 

```shell
make checksum                                                                                                          
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
```

Multi-tail the streams consumers, they should all be exactly the same (requires `multitail`)

```shell
make read-all
...
18 | block start | 0x
19 | transaction | 0x307832373434643137323732393864333838343862636165356630633539383062313337383862323266613163653433656336386664633733363030383766336435
20 |   block end | 0x
21 | block start | 0x
22 | transaction | 0x307862656532626666613738303939393766323263383664356663316161646664386536376638343733376364626663643262383332336666353938313739633931
23 | transaction | 0x307832306438393862393563303164323838313766633164643962366462363835363266626535623232356463643131663732326638353266623033326231333963
24 | transaction | 0x307833616665633731383738383464633061346266383366373765623264623832383732373039653337663661666465656461653530316263613466363063353334
25 |   block end | 0x
```
