package cobra

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (cli *cliAdapter) vaultAddCredsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "creds",
		Short: "Add new creds to the vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			login, _ := cmd.Flags().GetString("login")
			password, _ := cmd.Flags().GetString("password")
			info, _ := cmd.Flags().GetString("info")

			data := domain.VaultDataCreds{
				Login:    login,
				Password: password,
				Info:     info,
			}

			_, err := cli.clientService.VaultAddCreds(cmd.Context(), data)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("New creds added successfully")

			return nil
		},
	}

	c.Flags().StringP("login", "l", "", "Login")
	c.Flags().StringP("password", "p", "", "Password")
	c.Flags().StringP("info", "i", "", "Metadata")
	c.MarkFlagsOneRequired("login", "password")

	return c
}
