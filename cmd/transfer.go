package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"encoding/json"
	"github.com/primasio/contract-safe-deploy/contracts"
)

func TransferCmd() cli.Command {
	return cli.Command{
		Name:    "transfer",
		Aliases: []string{"tf"},
		Usage:   "transfer token",
		Flags: []cli.Flag{
			inputFileFlag(),
			outputFileFlag(),
		},
		Action: func(c *cli.Context) error {
			cfg.InitEthClient()
			contracts.Init()

			logger.Info("input file", "file", cfg.IFile)

			d, err := ioutil.ReadFile(cfg.IFile)
			if err != nil {
				panic(err.Error())
			}

			iFile := make([]string, 0)
			if err := json.Unmarshal(d, &iFile); err != nil {
				panic(err.Error())
			}

			oFile := make([]string, 0)
			for _, ifile := range iFile {
				logger.Info("handle rule", "rule", ifile)

				das := strings.Split(ifile, ",")
				if len(das) != 4 {
					logger.Error("参数不够", "params", ifile)
					panic("请传入四个参数[contratName,contratMethod,userAddress,roleType]")
				}

				contratName := das[0]
				contratMth := das[1]
				amount := das[2]
				address := das[3]

				oFile = append(oFile, common.ToHex(contracts.GetContract(contratName).Transfer(contratMth, amount, address)))
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
