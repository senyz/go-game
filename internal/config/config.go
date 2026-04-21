package config

import "github.com/spf13/viper"

type Config struct {
    Server struct {
        Port int    `mapstructure:"port"`
        Host string `mapstructure:"host"`
    } `mapstructure:"server"`
    Database struct {
        Dialect string `mapstructure:"dialect"`
        DSN     string `mapstructure:"dsn"`
    } `mapstructure:"database"`
    Messenger struct {
        APIURL string `mapstructure:"api_url"`
        Token  string `mapstructure:"token"`
    } `mapstructure:"messenger"`
    Logging struct {
        Level  string `mapstructure:"level"`
        Format string `mapstructure:"format"`
    } `mapstructure:"logging"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("configs/")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
