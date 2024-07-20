package cobra

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

// initConfig sets the config files location and attempts to read it in.
func (cli *cliAdapter) initConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := path.Join(homeDir, ".config")

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModeDir|0755)
		if err != nil {
			return err
		}
	}

	configFile := path.Join(configDir, ".gophkeeper")

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		err = os.WriteFile(configFile, []byte(""), 0600)
		if err != nil {
			return err
		}
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	viper.SetDefault("auth.token", "")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
