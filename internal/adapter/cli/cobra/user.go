package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the гser registration and authentication command
func (cli *cliAdapter) userCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "user",
		Short: "User registration and authentication",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(cli.userLoginCmd())
	c.AddCommand(cli.userRegisterCmd())

	return c
}

func (cli *cliAdapter) userLoginCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "login",
		Short: "User authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			login, _ := cmd.Flags().GetString("login")
			password, _ := cmd.Flags().GetString("password")

			token, err := cli.clientService.UserLogin(cmd.Context(), login, password)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Successfully logged in")

			viper.Set("auth.token", token)
			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("could not write config: %w", err)
			}

			return nil
		},
	}

	c.Flags().String("login", "", "User login")
	c.Flags().String("password", "", "User password")

	return c
}

func (cli *cliAdapter) userRegisterCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "register",
		Short: "User registration",
		RunE: func(cmd *cobra.Command, args []string) error {
			login, _ := cmd.Flags().GetString("login")
			password, _ := cmd.Flags().GetString("password")

			token, err := cli.clientService.UserRegister(cmd.Context(), login, password)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("Successfully registered")

			viper.Set("auth.token", token)
			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("could not write config: %w", err)
			}

			return nil
		},
	}

	c.Flags().String("login", "", "User login")
	c.Flags().String("password", "", "User password")

	return c
}
