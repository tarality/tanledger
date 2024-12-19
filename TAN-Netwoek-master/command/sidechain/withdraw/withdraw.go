package withdraw

import (
	"fmt"
	"time"

	"github.com/tarality/tan-network/command"
	"github.com/tarality/tan-network/command/bridge/common"
	"github.com/tarality/tan-network/command/helper"
	"github.com/tarality/tan-network/command/polybftsecrets"
	sidechainHelper "github.com/tarality/tan-network/command/sidechain"
	"github.com/tarality/tan-network/consensus/polybft/contractsapi"
	"github.com/tarality/tan-network/contracts"
	"github.com/tarality/tan-network/txrelayer"
	"github.com/tarality/tan-network/types"
	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
)

var params withdrawParams

func GetCommand() *cobra.Command {
	unstakeCmd := &cobra.Command{
		Use:     "withdraw-child",
		Short:   "Withdraws pending withdrawals on child chain for given validator",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	helper.RegisterJSONRPCFlag(unstakeCmd)
	setFlags(unstakeCmd)

	return unstakeCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		polybftsecrets.AccountDirFlag,
		"",
		polybftsecrets.AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.accountConfig,
		polybftsecrets.AccountConfigFlag,
		"",
		polybftsecrets.AccountConfigFlagDesc,
	)

	cmd.MarkFlagsMutuallyExclusive(polybftsecrets.AccountDirFlag, polybftsecrets.AccountConfigFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	return params.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) error {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	validatorAccount, err := sidechainHelper.GetAccount(params.accountDir, params.accountConfig)
	if err != nil {
		return err
	}

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC),
		txrelayer.WithReceiptTimeout(150*time.Millisecond))
	if err != nil {
		return err
	}

	encoded, err := contractsapi.ValidatorSet.Abi.Methods["withdraw"].Encode([]interface{}{})
	if err != nil {
		return err
	}

	txn := &ethgo.Transaction{
		From:  validatorAccount.Ecdsa.Address(),
		Input: encoded,
		To:    (*ethgo.Address)(&contracts.ValidatorSetContract),
	}

	receipt, err := txRelayer.SendTransaction(txn, validatorAccount.Ecdsa)
	if err != nil {
		return err
	}

	if receipt.Status != uint64(types.ReceiptSuccess) {
		return fmt.Errorf("withdraw transaction failed on block: %d", receipt.BlockNumber)
	}

	var (
		withdrawalEvent contractsapi.WithdrawalEvent
		foundLog        bool
	)

	// check the logs to check for the result
	for _, log := range receipt.Logs {
		doesMatch, err := withdrawalEvent.ParseLog(log)
		if err != nil {
			return err
		}

		if doesMatch {
			foundLog = true

			break
		}
	}

	if !foundLog {
		return fmt.Errorf("could not find an appropriate log in receipt that withdraw happened on ValidatorSet")
	}

	exitEventID, err := common.ExtractExitEventID(receipt)
	if err != nil {
		return fmt.Errorf("withdrawal failed: %w", err)
	}

	outputter.WriteCommandResult(
		&withdrawResult{
			ValidatorAddress: validatorAccount.Ecdsa.Address().String(),
			Amount:           withdrawalEvent.Amount,
			ExitEventID:      exitEventID,
			BlockNumber:      receipt.BlockNumber,
		})

	return nil
}
