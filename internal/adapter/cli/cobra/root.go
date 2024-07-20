package cobra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/util"
)

// rootCmd represents the base command when called without any subcommands
func (cli *cliAdapter) rootCmd(buildDate, buildVersion string) *cobra.Command {
	c := &cobra.Command{
		Use:   "gophkeeper",
		Short: "Application for storing confidential information",
		Long: func() string {
			return "Application for storing confidential information.\n" +
				util.GetBuildInfo(buildDate, buildVersion)
		}(),
	}

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		token := viper.GetString("auth.token")
		cli.clientService.SetToken(token)
	}

	c.AddCommand(cli.userCmd())

	return c
}
