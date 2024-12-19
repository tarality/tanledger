package ibft

import (
	"github.com/tarality/tan-network/command/helper"
	"github.com/tarality/tan-network/command/ibft/candidates"
	"github.com/tarality/tan-network/command/ibft/propose"
	"github.com/tarality/tan-network/command/ibft/quorum"
	"github.com/tarality/tan-network/command/ibft/snapshot"
	"github.com/tarality/tan-network/command/ibft/status"
	_switch "github.com/tarality/tan-network/command/ibft/switch"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	ibftCmd := &cobra.Command{
		Use:   "ibft",
		Short: "Top level IBFT command for interacting with the IBFT consensus. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(ibftCmd)

	registerSubcommands(ibftCmd)

	return ibftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// ibft status
		status.GetCommand(),
		// ibft snapshot
		snapshot.GetCommand(),
		// ibft propose
		propose.GetCommand(),
		// ibft candidates
		candidates.GetCommand(),
		// ibft switch
		_switch.GetCommand(),
		// ibft quorum
		quorum.GetCommand(),
	)
}
