# [![golang](https://raw.githubusercontent.com/the-code-innovator/go-blockchain/master/assets/assets/mascot.png)](https://golang.org)**golang** Block Chain

[![MIT License](https://img.shields.io/cocoapods/l/AFNetworking.svg?style=plastic)](https://github.com/the-code-innovator/go-blockchain/blob/master/LICENSE) [![Join the chat at https://gitter.im/go-blockchain-/community](https://badges.gitter.im/go-blockchain-/community.svg)](https://gitter.im/go-blockchain-/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![GitHub](https://raw.githubusercontent.com/the-code-innovator/go-blockchain/master/assets/assets/github.png)](https://github.com/the-code-innovator/go-blockchain)
[![Gitter](https://raw.githubusercontent.com/the-code-innovator/go-blockchain/master/assets/assets/gitter.png)](https://gitter.im/go-blockchain-/community)

## [Blockchain](https://en.wikipedia.org/wiki/Blockchain)

* A blockchain, originally block chain, is a growing list of records, called blocks, which are linked using cryptography.
* Each block contains a cryptographic hash of the previous block, a timestamp, and transaction data (generally represented as a merkle tree root hash).

## Description

a block chain in [**golang**](https://golang.org) with command line interface builtin.

## Language

* [golang](https://golang.org)

## Usage

```text
USAGE:
    -> dev   : go run main.go   <OPTIONS>
    -> build : go-blockchain <OPTIONS>
 • getbalance -address ADDRESS           - get balance for address.
 • createblockchain -address ADDRESS     - creates a blockchain.
 • printchain                            - prints the blocks in the blockchain.
 • send -from FROM -to TO -amount AMOUNT - send amount from an address to an address.
 • createwallet                          - creates a new wallet.
 • listaddresses                         - lists the addresses in our wallet file.
```

## Utilities

* getbalance:
   ```$ $EXECUTABLE getbalance -address ADDRESS```
  * To get balance of address 'ADDRESS'.
* createblockchain:
   ```$ $EXECUTABLE createblockchain -address ADDRESS```
  * To create a blockchain and send reward to the address 'ADDRESS'.
* printchain:
   ```$ $EXECUTABLE printchain```
  * To print the blocks in the blockchain.
* send:
   ```$ $EXECUTABLE send -from FROM -to TO -amount AMOUNT```
  * To send amount AMOUNT from address 'FROM' to address 'TO'.
* createwallet:
   ```$ $EXECUTABLE createwallet```
  * To create a wallet and store it in the wallets database.
* listaddresses:
   ```$ $EXECUTABLE listaddresses```
  * To list all public addresses in the wallets database.

`$EXECUTABLE` evaluvates to:

* dev:
      ```go run main.go```
* build:
  * First build the module.
      ```$ go build```
  * After building the module.
      ```$PWD/go-blockchain```
