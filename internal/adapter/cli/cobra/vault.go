package cobra

import "github.com/spf13/cobra"

func (cli *cliAdapter) vaultCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "vault",
		Short: "Actions with the vault",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(cli.vaultAddCmd())
	c.AddCommand(cli.vaultListCmd())
	c.AddCommand(cli.vaultGetCmd())

	return c
}
