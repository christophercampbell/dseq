# desq

distributed sequencer built on cometbft

https://docs.cometbft.com/v0.37/guides/go-built-in

### build the binary

```shell
make build
make build-linux
```

Cometbft generates a test cluster of nodes that are configured as peers. The docker compose file defines 4 nodes that use these configurations.

To run the local testnet cluster

```shell
make start
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
$ make compare                                                                                                          
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
a1202e1cd87dde92e86e53c1c12905e2
```

Since the image mounts the local build directory, the binary can be rebuilt without having to rebuild the local docker image

```shell
make build-linux
```

Send a transaction

```
curl -s 'localhost:26657/broadcast_tx_commit?tx="cometbft=rocks"'
```
