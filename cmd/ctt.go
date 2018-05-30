package cmd

import (
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/hex"
	"github.com/kooksee/pstoff/contracts"
	"encoding/json"
)

func CTTCmd() cli.Command {
	return cli.Command{
		Name:    "contract",
		Aliases: []string{"ctt"},
		Usage:   "contract opt",
		Flags: []cli.Flag{
			inputFileFlag(),
			outputFileFlag(),
		},
		Action: func(c *cli.Context) error {
			d, err := ioutil.ReadFile(cfg.IFile)
			if err != nil {
				panic(err.Error())
			}

			iFile := make(map[string][]string)
			if err := json.Unmarshal(d, &iFile); err != nil {
				panic(err.Error())
			}

			oFile := make([]string, 0)
			for key, ifile := range iFile {
				data := strings.Split(key, ".")
				if len(data) < 2 {
					panic(fmt.Sprintf("该数据 ％s 解析失败", ifile))
				}

				contractName := data[0]
				contracrMethod := data[1]

				args1 := []interface{}{}
				for _, a := range ifile {
					data := strings.Split(a, ".")
					d := data[0]
					t := data[1]
					if t == "bytes" {
						d1, err := hex.DecodeString(d)
						if err != nil {
							panic(err.Error())
						}
						args1 = append(args1, d1)
					} else {
						args1 = append(args1, d)
					}
				}

				h := contracts.GetContract(contractName).Execute(contracrMethod, args1...)
				oFile = append(oFile, hex.EncodeToString(h))
			}

			d1, err := json.Marshal(oFile)
			if err != nil {
				panic(err.Error())
			}
			if err := ioutil.WriteFile(cfg.OFile, d1, 0755); err != nil {
				panic(fmt.Sprintf("写入失败\n%s", err.Error()))
			}

			return nil
		},
	}
}
