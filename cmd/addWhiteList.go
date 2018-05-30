package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"github.com/kooksee/pstoff/contracts"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"encoding/json"
)

func AddWhiteListCmd() cli.Command {
	return cli.Command{
		Name:    "whiteList",
		Aliases: []string{"wt"},
		Usage:   "add white list",
		Flags: []cli.Flag{
			inputFileFlag(),
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
				logger.Info("handle white list", "list", ifile)

				das := strings.Split(ifile, ",")
				if len(das) != 3 {
					logger.Error("参数不够", "params", ifile)
					panic("请传入三个参数[contratName,contratMethod,address]")
				}

				contratName := das[0]
				contratMth := das[1]
				userAddress := das[2]

				oFile = append(oFile, common.ToHex(contracts.GetContract(contratName).AddWhiteList(contratMth, userAddress)))
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
