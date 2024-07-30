package cobra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (cli *cliAdapter) vaultAddFileCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "file",
		Short: "Add a new file to the vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("file")
			info, _ := cmd.Flags().GetString("info")

			filename := filepath.Base(path)

			filebytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			mime, err := mimetype.DetectFile(path)
			if err != nil {
				return err
			}

			data := domain.VaultDataFile{
				Filename: filename,
				MimeType: mime.String(),
				Filesize: uint64(len(filebytes)),
				Data:     filebytes,
				Info:     info,
			}

			_, err = cli.clientService.VaultAddFile(cmd.Context(), data)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("New file added successfully")

			return nil
		},
	}

	c.Flags().StringP("file", "f", "", "Path to the file")
	c.Flags().StringP("info", "i", "", "File metadata")
	c.MarkFlagRequired("file")

	return c
}
