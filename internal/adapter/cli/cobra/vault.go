package cobra

import "github.com/spf13/cobra"

func (cli *cliAdapter) vaultCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "vault",
		Short: "Actions with the vault",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(cli.vaultAddCmd())

	return c
}

func (cli *cliAdapter) vaultAddCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Add a new record to the vault",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(cli.vaultAddNoteCmd())

	return c
}
