package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jackvonhouse/auth-service/pkg/log"
	"github.com/spf13/viper"
)

type Token struct {
	Exp int
}

type JWT struct {
	AccessToken  *Token
	RefreshToken *Token
	SecretKey    string
}

type Database struct {
	Username string
	Password string
	Cluster  string
}

func (d *Database) String() string {
	return fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		d.Username, d.Password, d.Cluster,
	)
}

type ServerHTTP struct {
	Port int
}

type Config struct {
	Database *Database
	JWT      *JWT
	Server   *ServerHTTP
}

func New(
	configPath string,
	logger log.Logger,
) (*Config, error) {

	configType := strings.TrimPrefix(filepath.Ext(configPath), ".")

	viper.SetConfigType(configType)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(map[string]any{
			"layer":       "config",
			"config_path": configPath,
		}).Warnf("error on reading config: %s", err)

		return nil, fmt.Errorf("error on reading config: %s", err)
	}

	mongoPrefix := "database.mongo"
	tokenPrefix := "token"

	return &Config{
		Database: &Database{
			Username: viper.GetString(
				fmt.Sprintf("%s.username", mongoPrefix),
			),

			Password: viper.GetString(
				fmt.Sprintf("%s.password", mongoPrefix),
			),

			Cluster: viper.GetString(
				fmt.Sprintf("%s.cluster", mongoPrefix),
			),
		},

		JWT: &JWT{
			AccessToken: &Token{
				Exp: viper.GetInt(
					fmt.Sprintf("%s.access.exp", tokenPrefix),
				),
			},

			RefreshToken: &Token{
				Exp: viper.GetInt(
					fmt.Sprintf("%s.refresh.exp", tokenPrefix),
				),
			},

			SecretKey: viper.GetString(
				fmt.Sprintf("%s.secret", tokenPrefix),
			),
		},

		Server: &ServerHTTP{
			Port: viper.GetInt("server.http.port"),
		},
	}, nil
}
