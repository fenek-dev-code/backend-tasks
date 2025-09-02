package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DataBaseConfig `env:"storage" yaml:"storage" env-default:""`
	SqliteConfig   `env:"sqlite" yaml:"sqlite" env-default:""`
	ServerConfig   `env:"server" yaml:"server" env-default:""`
	LoggerConfig   `env:"logger" yaml:"logger" env-default:""`
}

type LoggerConfig struct {
	LogLevl string `env:"LOG_LEVEL,required" yaml:"log_level" env-default:"debug"`
	LogPath string `env:"LOG_PATH,required" yaml:"log_path" env-default:"./logs/log.log"`
}

type DataBaseConfig struct {
	PostgresHost     string `env:"POSTGRES_HOST,required" yaml:"POSTGRES_HOST" env-default:"localhost" `
	PostgresPort     int    `env:"POSTGRES_PORT,required" yaml:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER,required" yamal:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required" yaml:"POSTGRES_PASSWORD" env-default:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB,required" yaml:"POSTGRES_DB" env-default:"postgres"`
}

type ServerConfig struct {
	Prefix string `env:"PREFIX,required" yaml:"PREFIX" env-default:"/api/v1"`
	Host   string `env:"HOST,required" yaml:"HOST" env-default:"localhost"`
	Port   int    `env:"PORT,required" yaml:"PORT" env-default:"8080"`
}

type SqliteConfig struct {
	Path string `env:"SQLITE_PATH,required" yaml:"path" env-default:"./tasks.db"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("config.yml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *DataBaseConfig) GetPostgresURL() string {
	url := "posrgres://"
	url += c.PostgresUser + ":" + c.PostgresPassword
	url += "@" + c.PostgresHost + ":" + string(rune(c.PostgresPort))
	return url
}
