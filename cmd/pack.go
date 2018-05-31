package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/primasio/contract-safe-deploy/contracts"
)

func PackCmd() cli.Command {
	return cli.Command{
		Name:    "pack",
		Aliases: []string{"pk"},
		Usage:   "pack method",
		Flags: []cli.Flag{
			inputFileFlag(),
		},
		Action: func(c *cli.Context) error {
			contracts.Init()

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
				das := strings.Split(t, ",")
				contactName := das[0]

				logger.Info("contract name","name",contactName)
				logger.Info("pack data","data",t)
				ct := contracts.GetContract(contactName)

				var adds []interface{}
				for _, i := range das[1:] {
					adds = append(adds, common.HexToAddress(i))
				}
				
				dd, err := ct.ABI.Pack("", adds...)
				if err != nil {
					panic(err.Error())
				}

				dd = append(ct.Btc, dd...)
				oFile = append(oFile, common.ToHex(dd))
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
