package main

import (
	"github.com/primasio/contract-safe-deploy/config"
	"github.com/primasio/contract-safe-deploy/cmd"
)

func main() {
	cfg := config.NewCfg("kdata")
	cfg.InitLog()
	cfg.LoadConfig()
	cfg.InitNode()

	cmd.Init()
	cmd.RunCmd()
}
