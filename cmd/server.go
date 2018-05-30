package cmd

import (
	"github.com/urfave/cli"
)

func ServerCmd() cli.Command {
	return cli.Command{
		Action: func(c *cli.Context) error {

			//node.New(cfg.Seeds).RunHttpServer()

			return nil
		},
	}
}
