package cobra

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (cli *cliAdapter) vaultGetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Get vault record or save file by id",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(cli.vaultGetRecordCmd())
	c.AddCommand(cli.vaultGetFileCmd())

	return c
}

func (cli *cliAdapter) vaultGetRecordCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "record <id>",
		Short: "Get vault record by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recordID, err := uuid.Parse(args[0])
			if err != nil {
				return err
			}

			record, err := cli.clientService.VaultGetRecord(cmd.Context(), recordID)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			props := map[string]string{}
			props["Id"] = record.ID.String()
			props["Kind"] = record.Kind.String()

			switch record.Kind {
			case domain.VaultKindNote:
				data := record.Data.(domain.VaultDataNote)
				props["Info"] = data.Info
				props["Content"] = data.Content
			case domain.VaultKindCard:
				data := record.Data.(domain.VaultDataCard)
				props["Info"] = data.Info
				props["Number"] = data.Number
				props["Holder Name"] = data.HolderName
				props["Exp. date"] = fmt.Sprintf("%02d/%02d", data.ExpDate.Month, data.ExpDate.Year)
				props["CVC code"] = data.CvcCode
			case domain.VaultKindCreds:
				data := record.Data.(domain.VaultDataCreds)
				props["Info"] = data.Info
				props["Login"] = data.Login
				props["Password"] = data.Password
			case domain.VaultKindFile:
				data := record.Data.(domain.VaultDataFile)
				props["Info"] = data.Info
				props["Filename"] = data.Filename
				props["MimeType"] = data.MimeType
				props["Filesize"] = strconv.FormatUint(data.Filesize, 10)
			}

			props["Created at"] = record.CreatedAt.Format("2006-01-02 15:04:05")
			props["Updated at"] = record.CreatedAt.Format("2006-01-02 15:04:05")

			fmt.Printf("%-20s | %-80s |\n", "Property", "Value")
			fmt.Printf("%-20s | %-80s |\n", strings.Repeat("_", 20), strings.Repeat("_", 80))
			for k, v := range props {
				fmt.Printf("%-20s | %-80s |\n", k, v)
			}
			fmt.Printf("%-20s | %-80s |\n", strings.Repeat("_", 20), strings.Repeat("_", 80))

			return nil
		},
	}

	c.Flags().Uint64P("page", "p", 1, "Page number")
	c.Flags().Uint64P("limit", "l", 10, "Number of records per page")

	return c
}

func (cli *cliAdapter) vaultGetFileCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "file <id>",
		Short: "Download file from vault by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recordID, err := uuid.Parse(args[0])
			if err != nil {
				return err
			}

			file, data, err := cli.clientService.VaultGetFile(cmd.Context(), recordID)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			ex, err := os.Executable()
			if err != nil {
				return err
			}
			path := filepath.Dir(ex)

			filename := path + "/" + file.Filename

			err = os.WriteFile(filename, data, 0644)
			if err != nil {
				fmt.Println("Failed to save file:", err)
				return nil
			}

			fmt.Println("File saved successfully to:", filename)

			return nil
		},
	}

	return c
}
