package contracts

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/primasio/contract-safe-deploy/config"
)

var cfg *config.Config
var logger log15.Logger

func Init() {
	cfg = config.GetCfg()
	logger = config.Log().New("package", "contracts")
	if err := initContracts(); err != nil {
		panic(fmt.Sprintf("初始化合约错误\n%s", err.Error()))
	}
}
