package configs

import (
	"os"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const RequestUserID = "userID"
const RequestSID = "SID"

// type DatabaseConfig struct {
// 	Host            string `yaml:"host"`
// 	Port            int    `yaml:"port"`
// 	Database        string `yaml:"dbname"`
// 	User            string `yaml:"user"`
// 	Password        string `yaml:"password"`
// 	SetMaxOpenConns int    `yaml:"max_open_conns"`
// }

type Config struct {
	ApiPath    string           `yaml:"apiPath"`
	Server     ServerConfig     `json:"server" yaml:"server"`
	Database   DatabaseConfig   `json:"database" yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	FilesPaths FilesPathsConfig `json:"filesPaths" yaml:"filesPaths"`
}

type ServerConfig struct {
	Host        string        `json:"host" yaml:"host"`
	Port        string        `json:"port" yaml:"port"`
	SwaggerPort string        `json:"swaggerPort" yaml:"swaggerPort"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	//Yoomoney    Yoomoney      `yaml:"yoomoney"`
}

type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Database        string `json:"database"`
	User            string `json:"username"`
	Password        string `json:"password"`
	Timer           uint32 `json:"timer"`
	GrpcPort        string `json:"grpc_port"`
	SetMaxOpenConns int    `json:"set_max_open_conns"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AwsConfig struct {
	Id     string `json:"key_id" yaml:"id"`
	Access string `json:"key_access" yaml:"access"`
	Region string `json:"region" yaml:"region"`
}

type FilesPathsConfig struct {
}

var Cfg Config

func ReadConfig() (*DatabaseConfig, error) {
	dsnConfig := DatabaseConfig{}
	dsnFile, err := os.ReadFile("../../configs/db_dsn.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dsnFile, &dsnConfig)
	if err != nil {
		return nil, err
	}

	return &dsnConfig, nil

}

func LoadConfig(path string) (*Config, error) {
	var err error
	var config Config
	viper.SetConfigFile(path)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err = viper.BindEnv("server.host", "SERVER_HOST")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("server.port", "SERVER_PORT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("server.timeout", "SERVER_TIMEOUT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.host", "DB_HOST")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.port", "DB_PORT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.user", "DB_USER")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.password", "DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.dbname", "DB_NAME")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.timer", "DB_TIMER")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.grpc", "DB_GRPC")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("aws.key_id", "AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("aws.key_access", "AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("aws.region", "AWS_DEFAULT_REGION")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("yoomoney.key", "YOOMONEY_KEY")
	if err != nil {
		return nil, err
	}
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	Cfg = config
	return &config, nil
}
