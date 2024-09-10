package util

import "github.com/spf13/viper"

type Config struct {
	EthClientAPIKey    string `mapstructure:"ETH_API_KEY"`
	UniethTokenAddress string `mapstructure:"UNIETH_TOKEN_ADDRESS"`
	UserAddress        string `mapstructure:"USER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
