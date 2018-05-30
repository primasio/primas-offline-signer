package cmd

import (
	"github.com/urfave/cli"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"encoding/json"
)

const abiStr = `[
	{
		"constant": true,
		"inputs": [],
		"name": "get",
		"outputs": [
			{
				"name": "retVal",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "x",
				"type": "string"
			}
		],
		"name": "set",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

const contractAddr = "0xC1324b9c86B5FD6C3C6fA08738DBf2909fb27C16"

func TestCmd() cli.Command {
	return cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "test",
		Action: func(c *cli.Context) error {
			cfg.InitEthClient()

			abiApi, err := abi.JSON(strings.NewReader(abiStr))
			if err != nil {
				panic(err.Error())
			}

			methodBytes, err := abiApi.Pack("set", "hello")
			if err != nil {
				panic(err.Error())
			}

			tx := types.NewTransaction(
				cfg.GetNonce(),
				common.HexToAddress(contractAddr),
				big.NewInt(0),
				big.NewInt(int64(cfg.GasLimit)),
				big.NewInt(int64(cfg.Gasprice)),
				methodBytes,
			)

			dd1, err := json.Marshal(tx)
			if err != nil {
				panic(err.Error())
			}

			logger.Error(string(dd1))

			logger.Info(cfg.GetNodeAccount().Address.Hex())

			ks := cfg.GetNodeKeyStore()
			signedTx, err := ks.SignTx(*cfg.GetNodeAccount(), tx, big.NewInt(int64(3)))
			if err != nil {
				panic(err.Error())
			}

			dd, err := signedTx.MarshalJSON()
			if err != nil {
				panic(err.Error())
			}

			logger.Error(string(dd))

			//client := cfg.GetEthClient()
			//ctx2, _ := context.WithTimeout(context.Background(), time.Minute)
			//if err := client.SendTransaction(ctx2, signedTx); err != nil {
			//	panic(err.Error())
			//}

			logger.Info("SendTransaction", "tx hash", signedTx.Hash().String())
			return nil
		},
	}
}

//tx1 := &types.Transaction{}
//if err := tx1.UnmarshalJSON([]byte(`{"nonce":"0xc9c6","gasPrice":"0x174876e800","gas":"0xb71b0","to":"0xc1324b9c86b5fd6c3c6fa08738dbf2909fb27c16","value":"0x0","input":"0x4ed3885e0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000568656c6c6f000000000000000000000000000000000000000000000000000000","v":"0x2a","r":"0x37ebb3af377e08709fe83f41b8767fee908ea8a21ae6fd40a1c552abd89f4d3a","s":"0x2aaf38d4c4372aabc9ce11208aca3172a5cff0d5ea1d318ee5c44a6388bed76f","hash":"0xaa660f3da1e3bfe9a8921701263ab9017c6eb918f53e4ab5709dd391ce4ef066"}`)); err != nil {
//	panic(err.Error())
//}

//d, _ := tx1.MarshalJSON()
//logger.Error(string(d))
