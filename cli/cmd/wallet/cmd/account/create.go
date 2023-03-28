package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// createCmd represents the create account command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a wallet account.",
	Long:  `This command creates an account using seed/validator private key and index (index_from).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Create(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddIndexFlag(createCmd)
	flag.AddSeedFlag(createCmd)
	flag.AddPrivateKeyFlag(createCmd)
	rootcmd.AddAccumulateFlag(createCmd)
	flag.AddResponseTypeFlag(createCmd)
	flag.AddHighestSourceFlag(createCmd)
	flag.AddHighestTargetFlag(createCmd)
	flag.AddHighestProposalFlag(createCmd)
	rootcmd.AddNetworkFlag(createCmd)

	Command.AddCommand(createCmd)
}
