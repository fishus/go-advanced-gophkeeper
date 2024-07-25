package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cli *cliAdapter) vaultListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List of records from vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			page, _ := cmd.Flags().GetUint64("page")
			limit, _ := cmd.Flags().GetUint64("limit")

			list, err := cli.clientService.VaultListRecords(cmd.Context(), page, limit)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Printf("%-36s | %-5s | %-40s | %-20s | %-20s\n", "UUID", "KIND", "INFO", "CREATED", "UPDATED")
			for _, item := range list {
				fmt.Printf("%-36s | %-5s | %-40s | %-20s | %-20s\n", item.ID, item.Kind, item.Info, item.CreatedAt.Format("2006-01-02 15:04:05"), item.UpdatedAt.Format("2006-01-02 15:04:05"))
			}

			return nil
		},
	}

	c.Flags().Uint64P("page", "p", 1, "Page number")
	c.Flags().Uint64P("limit", "l", 10, "Number of records per page")

	return c
}
