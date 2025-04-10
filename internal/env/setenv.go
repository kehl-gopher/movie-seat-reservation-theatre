package env

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

func SetEnv() *Config {

	c := &Config{}
	envs := LoadEnv()

	err := GetEnv(envs, c)
	if err != nil {
		panic(err)
	}

	return c
}

func GetEnv(envs map[string]interface{}, config *Config) error {

	err := mapstructure.Decode(envs, config)

	if err != nil {
		return err
	}
	return nil
}
func LoadEnv() map[string]interface{} {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			panic(err.Error())
		}
	}

	return viper.AllSettings()
}
