package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tarality/tan-network/command/backup"
	"github.com/tarality/tan-network/command/bridge"
	"github.com/tarality/tan-network/command/genesis"
	"github.com/tarality/tan-network/command/helper"
	"github.com/tarality/tan-network/command/ibft"
	"github.com/tarality/tan-network/command/license"
	"github.com/tarality/tan-network/command/monitor"
	"github.com/tarality/tan-network/command/peers"
	"github.com/tarality/tan-network/command/polybft"
	"github.com/tarality/tan-network/command/polybftsecrets"
	"github.com/tarality/tan-network/command/regenesis"
	"github.com/tarality/tan-network/command/rootchain"
	"github.com/tarality/tan-network/command/secrets"
	"github.com/tarality/tan-network/command/server"
	"github.com/tarality/tan-network/command/status"
	"github.com/tarality/tan-network/command/txpool"
	"github.com/tarality/tan-network/command/version"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "TAN Network is a EVM Halving based Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		rootchain.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		license.GetCommand(),
		polybftsecrets.GetCommand(),
		polybft.GetCommand(),
		bridge.GetCommand(),
		regenesis.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
