package examples

import (
	"fmt"

	"github.com/thecatshidog/ethgo"
	"github.com/thecatshidog/ethgo/abi"
	"github.com/thecatshidog/ethgo/contract"
	"github.com/thecatshidog/ethgo/jsonrpc"
	"github.com/thecatshidog/ethgo/wallet"
)

func contractTransaction() {
	var functions = []string{
		"function transferFrom(address from, address to, uint256 value)",
	}

	abiContract, err := abi.NewABIFromList(functions)
	handleErr(err)

	// Matic token
	addr := ethgo.HexToAddress("0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0")

	client, err := jsonrpc.NewClient("https://mainnet.infura.io")
	handleErr(err)

	// wallet signer
	key, err := wallet.GenerateKey()
	handleErr(err)

	opts := []contract.ContractOption{
		contract.WithJsonRPC(client.Eth()),
		contract.WithSender(key),
	}
	c := contract.NewContract(addr, abiContract, opts...)
	txn, err := c.Txn("transferFrom", ethgo.Latest)
	handleErr(err)

	err = txn.Do()
	handleErr(err)

	receipt, err := txn.Wait()
	handleErr(err)

	fmt.Printf("Transaction mined at: %s", receipt.TransactionHash)
}
