package cmd

import (
	"log"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/metis/printing"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/rhea/deployments"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

var (
	SOURCE      = deployments.Prod_OptimismToGoerli
	DESTINATION = deployments.Prod_GoerliToOptimism
)

type MetisClient struct {
	Logger      logger.Logger
	CloseLogger func() error
}

func NewMetisApp(client MetisClient) *cli.App {
	app := cli.NewApp()
	app.Name = "Metis"
	app.Usage = "CCIP sanity checker"

	err := SOURCE.SetupReadOnlyChain(client.Logger.Named(helpers.ChainName(int64(SOURCE.ChainConfig.ChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = DESTINATION.SetupReadOnlyChain(client.Logger.Named(helpers.ChainName(int64(DESTINATION.ChainConfig.ChainId))))
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "state",
			Aliases: []string{"s"},
			Usage:   "prints CCIP config state",
			Action: func(c *cli.Context) error {
				printing.PrintCCIPState(&SOURCE, &DESTINATION)
				return nil
			},
		},
		{
			Name:    "txs",
			Aliases: []string{"t"},
			Usage:   "prints recent txs",
			Action: func(c *cli.Context) error {
				printing.PrintTxStatuses(&SOURCE, &DESTINATION)
				return nil
			},
		},
	}

	return app
}
