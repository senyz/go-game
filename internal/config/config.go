// internal/config/config.go
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
		APIURL     string `mapstructure:"api_url"`
		Token      string `mapstructure:"token"`
		WebhookURL string `mapstructure:"webhook_url"` // добавляем
		UseWebhook bool   `mapstructure:"use_webhook"` // добавляем
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

	// Устанавливаем значения по умолчанию
	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Messenger.WebhookURL == "" {
		cfg.Messenger.UseWebhook = false // если webhook URL не задан, используем long polling
	}

	return &cfg, nil
}
