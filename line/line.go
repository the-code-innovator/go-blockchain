package line

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/the-code-innovator/go-blockchain/blockchain"
	"github.com/the-code-innovator/go-blockchain/wallet"
)

// Interface struct for handling command line inter
type Interface struct{}

// PrintUsage for printing usage instructions
func (inter *Interface) PrintUsage() {
	fmt.Println("USAGE:")
	// fmt.Println("    -> dev   : go run main.go   <OPTIONS>")
	// fmt.Println("    -> build : ./go-block-chain <OPTIONS>")
	fmt.Println(" • getbalance -address ADDRESS           - get balance for address.")
	fmt.Println(" • createblockchain -address ADDRESS     - creates a blockchain.")
	fmt.Println(" • printchain                            - prints the blocks in the blockchain.")
	fmt.Println(" • send -from FROM -to TO -amount AMOUNT - send amount from an address to an address.")
	fmt.Println(" • createwallet                          - creates a new wallet.")
	fmt.Println(" • listaddresses                         - lists the addresses in our wallet file.")
}

// ValidateArguments to validate the arguments for the CommandInterface
func (inter *Interface) ValidateArguments() {
	if len(os.Args) < 2 {
		inter.PrintUsage()
		runtime.Goexit()
	}
}

// ListAddresses to list all addresses in the addressbook
func (inter *Interface) ListAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}

// CreateWallet to create a wallet in the addressbook
func (inter *Interface) CreateWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()
	fmt.Printf("NEW ADDRESS: %s\n", address)
}

// PrintChain to print the Blocks in the BlockChain from inter
func (inter *Interface) PrintChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.DataBase.Close()
	iterator := chain.Iterator()
	for {
		block := iterator.Next()
		fmt.Printf("PREVIOUS HASH: %x\n", block.PreviousHash)
		fmt.Printf("MAIN HASH: %x\n", block.Hash)
		proofOfWork := blockchain.NewProof(block)
		fmt.Printf("PROOF OF WORK: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

// CreateBlockChain to create a blockchain with the address as the genesis.
func (inter *Interface) CreateBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.DataBase.Close()
	fmt.Println("FINISHED CREATING BLOCKCHAIN.")
}

// GetBalance to get the balance from the address
func (inter *Interface) GetBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.DataBase.Close()
	balance := 0
	unSpentTransactionOutputs := chain.FindUnspentTransactionsOutputs(address)
	for _, output := range unSpentTransactionOutputs {
		balance += output.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// Send to send the amount from FROM to TO
func (inter *Interface) Send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.DataBase.Close()
	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("SUCCESS.")
}

// Run to run the inter
func (inter *Interface) Run() {
	inter.ValidateArguments()
	getBalanceCommand := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCommand := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCommand := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCommand := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCommand := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCommand.String("address", "", "The Address to find Balance.")
	createBlockChainAddress := createBlockChainCommand.String("address", "", "The Address to send Reward to.")
	sendFrom := sendCommand.String("from", "", "Source Wallet Address")
	sendTo := sendCommand.String("to", "", "Destination Wallet Address")
	sendAmount := sendCommand.Int("amount", 0, "Amount To Send")
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockChainCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "listaddresses":
		err := listAddressesCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		err := createWalletCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCommand.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		inter.PrintUsage()
		runtime.Goexit()
	}
	if getBalanceCommand.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCommand.Usage()
			runtime.Goexit()
		}
		inter.GetBalance(*getBalanceAddress)
	}
	if createBlockChainCommand.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCommand.Usage()
			runtime.Goexit()
		}
		inter.CreateBlockChain(*createBlockChainAddress)
	}
	if printChainCommand.Parsed() {
		inter.PrintChain()
	}
	if createWalletCommand.Parsed() {
		inter.CreateWallet()
	}
	if listAddressesCommand.Parsed() {
		inter.ListAddresses()
	}
	if sendCommand.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCommand.Usage()
			runtime.Goexit()
		}
		inter.Send(*sendFrom, *sendTo, *sendAmount)
	}
}
