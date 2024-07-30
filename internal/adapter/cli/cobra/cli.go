package cobra

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type cliAdapter struct {
	clientService port.ClientService
}

func New(clientService port.ClientService) *cliAdapter {
	adapter := cliAdapter{
		clientService: clientService,
	}
	return &adapter
}

// Execute adds all child commands to the root command and sets flags appropriately.
func (cli *cliAdapter) Execute(ctx context.Context, buildDate, buildVersion string) error {
	cobra.OnInitialize(func() {
		err := cli.initConfig()
		cobra.CheckErr(err)

		err = cli.clientService.Setup(ctx)
		cobra.CheckErr(err)
	})

	cobra.OnFinalize(func() {
		cli.clientService.Teardown(ctx)
	})

	err := cli.rootCmd(buildDate, buildVersion).ExecuteContext(ctx)
	return err
}
