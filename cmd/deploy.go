package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"encoding/json"
)

func DeployCmd() cli.Command {
	return cli.Command{
		Name:    "deploy",
		Aliases: []string{"dp"},
		Usage:   "deploy contract",
		Flags: []cli.Flag{
			inputFileFlag(),
			//outputFileFlag(),
		},
		Action: func(c *cli.Context) error {
			cfg.InitEthClient()
			//contracts.Init()

			d, err := ioutil.ReadFile(cfg.IFile)
			if err != nil {
				logger.Error("读文件失败", "err", err)
				panic(err.Error())
			}

			iFile := make([]string, 0)
			if err := json.Unmarshal(d, &iFile); err != nil {
				logger.Error("[]string json数据decode失败", "err", err)
				panic(err.Error())
			}

			oFile := make([]string, 0)
			for _, ifile := range iFile {
				d1 := common.FromHex(ifile)
				if d1 == nil {
					logger.Error("hex string decode error", "str", ifile)
					panic("")
				}
				oFile = append(oFile, common.ToHex(Deploy(d1)))
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
