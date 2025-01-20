package config

import (
	"github.com/spf13/viper"
)

// Config конфигурация клиента
type Config struct {
	v     *viper.Viper
	value *Value
}

// VK структура конфигурации
type VK struct {
	ApiVersion  string `mapstructure:"ApiVersion"`
	AccessToken string `mapstructure:"AccessToken"`
	GroupID     string `mapstructure:"GroupID"`
	AlbumID     string `mapstructure:"AlbumID"`
}

type Log struct {
	Level   string `mapstructure:"Level"`
	FlePath string `mapstructure:"FilePath"`
}

type Value struct {
	VK  VK  `mapstructure:"VK"`
	Log Log `mapstructure:"Log"`
}

// ErrorCfg ошибка конфигурации
type ErrorCfg error

// NewConfig конструктор
func NewConfig() (*Config, error) {
	var err error
	instance := new(Config)
	instance.v = viper.New()
	instance.value = new(Value)

	err = instance.init()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Config) init() error {
	var err error
	c.v.AddConfigPath(".")
	c.v.SetConfigName("client")
	c.v.SetConfigType("yaml")
	err = c.v.ReadInConfig()
	if err != nil {
		return ErrorCfg(err)
	}

	err = c.v.Unmarshal(c.value)
	if err != nil {
		return ErrorCfg(err)
	}

	return nil
}

func (c *Config) Value() *Value {
	return c.value
}
