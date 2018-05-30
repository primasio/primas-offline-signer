package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/core/types"
	kts "github.com/kooksee/pstoff/types"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"encoding/json"
	"math/big"
)

func SignCmd() cli.Command {
	return cli.Command{
		Name:    "sign",
		Aliases: []string{"sn"},
		Usage:   "sign tx",
		Flags: []cli.Flag{
			inputFileFlag(),
		},
		Action: func(c *cli.Context) error {

			logger.Info("input file", "file", cfg.IFile)

			d, err := ioutil.ReadFile(cfg.IFile)
			if err != nil {
				panic(err.Error())
			}

			ds := make([]string, 0)
			if err := json.Unmarshal(d, &ds); err != nil {
				panic(err.Error())
			}

			oFile := make([]string, 0)
			for _, t := range ds {
				tx := &kts.Tx{}
				tt1 := common.FromHex(t)
				if tt1 == nil {
					logger.Error("hex string error")
					panic("")
				}

				if err := tx.Decode(common.FromHex(t)); err != nil {
					logger.Error("decode tx error", "err", err)
					panic(err.Error())
				}

				logger.Info("info","data",string(tt1))

				var tx1 *types.Transaction
				if tx.IsCreateContract {
					tx1 = types.NewContractCreation(
						tx.Nonce,
						big.NewInt(tx.Amount),
						big.NewInt(tx.GasLimit),
						big.NewInt(tx.GasPrice),
						tx.Data,
					)
				} else {
					tx1 = types.NewTransaction(
						tx.Nonce,
						common.HexToAddress(tx.To),
						big.NewInt(tx.Amount),
						big.NewInt(tx.GasLimit),
						big.NewInt(tx.GasPrice),
						tx.Data,
					)
				}

				signedTx, err := cfg.GetNodeKeyStore().SignTx(*cfg.GetNodeAccount(), tx1, big.NewInt(int64(cfg.ChainId)))
				if err != nil {
					logger.Error("SignTx error", "err", err)
					panic(err.Error())
				}

				ddd, err := signedTx.MarshalJSON()
				if err != nil {
					panic(err.Error())
				}

				oFile = append(oFile, common.ToHex(ddd))
			}

			d1, err := json.Marshal(oFile)
			if err != nil {
				panic(err.Error())
			}

			cfg.OFile = cfg.IFile + fmt.Sprintf(".output.%d.json", cfg.Nonce)

			logger.Info("output file", "file", cfg.OFile)

			if err := ioutil.WriteFile(cfg.OFile, d1, 0755); err != nil {
				panic(fmt.Sprintf("写入失败\n%s", err.Error()))
			}

			return nil
		},
	}
}
