package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQHost      string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort      string `mapstructure:"RABBITMQ_PORT"`
}

func LoadConfig(paths ...string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

func (c *conf) String() string {
	return fmt.Sprintf(`
	DB: [Driver: %s, Host: %s, Port: %s, User: %s]
	RABBIT: [Host: %s, Port: %s]
	WEB: [Port: %s]
	GRPC: [Port: %s]
	GQL: [Port: %s]`,
		c.DBDriver, c.DBHost, c.DBPort, c.DBUser,
		c.RabbitMQHost, c.RabbitMQPort,
		c.WebServerPort, c.GRPCServerPort, c.GraphQLServerPort)
}
