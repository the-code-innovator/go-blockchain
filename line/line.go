package line

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/the-code-innovator/go-blockchain/blockchain"
	"github.com/the-code-innovator/go-blockchain/wallet"
)

const (
	major = 0
	minor = 1
	patch = 1
)

// Interface struct for handling command line inter
type Interface struct{}

// PrintInfo to print information of the system
func (inter *Interface) PrintInfo() {
	fmt.Printf("Go BlockChain - v%d.%d.%d\n", major, minor, patch)
}

// PrintUsage for printing usage instructions
func (inter *Interface) PrintUsage() {
	inter.PrintInfo()
	fmt.Println("USAGE:")
	fmt.Println(" • getbalance -address ADDRESS           - get balance for address.")
	fmt.Println(" • createblockchain -address ADDRESS     - creates a blockchain.")
	fmt.Println(" • printchain                            - prints the blocks in the blockchain.")
	fmt.Println(" • send -from FROM -to TO -amount AMOUNT - send amount from an address to an address.")
	fmt.Println(" • createwallet                          - creates a new wallet.")
	fmt.Println(" • listaddresses                         - lists the addresses in our wallet file.")
	fmt.Println(" • help                                  - prints the usage for the blockchain utility.")
}

// Help to print help information for the CommandInterface
func (inter *Interface) Help() {
	inter.PrintUsage()
	runtime.Goexit()
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
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()
		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

// CreateBlockChain to create a blockchain with the address as the genesis.
func (inter *Interface) CreateBlockChain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: ADDRESS IS NOT VALID !")
	}
	chain := blockchain.InitBlockChain(address)
	chain.DataBase.Close()
	fmt.Println("FINISHED CREATING BLOCKCHAIN.")
}

// GetBalance to get the balance from the address
func (inter *Interface) GetBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: ADDRESS IS NOT VALID !")
	}
	chain := blockchain.ContinueBlockChain(address)
	defer chain.DataBase.Close()
	balance := 0
	publicKeyHash := wallet.Base58Decode([]byte(address))
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	unSpentTransactionOutputs := chain.FindUnspentTransactionsOutputs(publicKeyHash)
	for _, output := range unSpentTransactionOutputs {
		balance += output.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// Send to send the amount from FROM to TO
func (inter *Interface) Send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: ADDRESS IS NOT VALID !")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: ADDRESS IS NOT VALID !")
	}
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
	case "help":
		inter.Help()
	default:
		inter.PrintInfo()
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
