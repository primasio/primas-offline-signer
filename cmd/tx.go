package cmd

import (
	"github.com/urfave/cli"
	"context"
	"time"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"encoding/json"
)

func TxCmd() cli.Command {
	return cli.Command{
		Name:    "sentTx",
		Aliases: []string{"st"},
		Usage:   "sent tx to chain",
		Flags: []cli.Flag{
			inputFileFlag(),
		},
		Action: func(c *cli.Context) error {
			cfg.InitEthClient()

			logger.Info("input file", "file", cfg.IFile)

			client := cfg.GetEthClient()
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
				tx := &types.Transaction{}
				tt1 := common.FromHex(t)
				if tt1 == nil {
					logger.Error("hex string error")
					panic("")
				}

				if err := tx.UnmarshalJSON(common.FromHex(t)); err != nil {
					logger.Error("decode tx error", "err", err)
					panic(err.Error())
				}

				ctx2, _ := context.WithTimeout(context.Background(), time.Minute)
				if err := client.SendTransaction(ctx2, tx); err != nil {
					logger.Error("SendTransaction  error", "err", err)
					panic(err.Error())
				}

				oFile = append(oFile, tx.Hash().String())
				logger.Info("SendTransaction", "hash", tx.Hash().String())
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
