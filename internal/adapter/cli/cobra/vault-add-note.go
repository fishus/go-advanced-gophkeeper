package cobra

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (cli *cliAdapter) vaultAddNoteCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "note",
		Short: "Add a new note to the vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			content, _ := cmd.Flags().GetString("content")
			info, _ := cmd.Flags().GetString("info")

			data := domain.VaultDataNote{
				Content: content,
				Info:    info,
			}

			rec, err := cli.clientService.VaultAddNote(cmd.Context(), data)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if rec != nil {
				fmt.Println("New note added successfully")
			} else {
				err = domain.ErrVaultRecordNotCreated
			}

			return nil
		},
	}

	c.Flags().StringP("content", "c", "", "Note content")
	c.Flags().StringP("info", "i", "", "Note metadata")
	_ = c.MarkFlagRequired("content")

	return c
}
