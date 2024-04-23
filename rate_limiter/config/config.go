package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	TokenLimit     int           `mapstructure:"TOKEN_LIMIT"`
	TokenBlockTime time.Duration `mapstructure:"TOKEN_BLOCK_TIME"`
	TokenFrameTime time.Duration `mapstructure:"TOKEN_FRAME_TIME"`
	IpLimit        int           `mapstructure:"IP_LIMIT"`
	IpBlockTime    time.Duration `mapstructure:"IP_BLOCK_TIME"`
	IpFrameTime    time.Duration `mapstructure:"IP_FRAME_TIME"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	Port int `mapstructure:"PORT"`
}

func (c Config) String() string {
	return fmt.Sprintf(`
TokenLimit: %d;
TokenBlockTime: %dms;
TokenFrameTime: %dms;
IpLimit: %d;
IpBlockTime: %dms;
IpFrameTime: %dms;
RedisAddr: %s;
RedisUsername: %s;
RedisPassword?: %t;
Port: %d;`,
		c.TokenLimit, c.TokenBlockTime.Milliseconds(), c.TokenFrameTime.Milliseconds(),
		c.IpLimit, c.IpBlockTime.Milliseconds(), c.IpFrameTime.Milliseconds(),
		c.RedisAddr, c.RedisUsername, c.RedisPassword != "",
		c.Port)
}

func LoadConfig(path string) *Config {
	var config *Config
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
