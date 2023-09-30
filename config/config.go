package config

type Config struct {
	Database struct {
		DBName   string `env:"DATABASE_NAME" env-default:"alireza"`
		Host     string `env:"DATABASE_HOST" env-default:"localhost"`
		Port     int    `env:"DATABASE_PORT" env-default:"5432"`
		Username string `env:"DATABASE_USER" env-default:"alireza"`
		Password string `env:"DATABASE_PASSWORD" env-default:"1234"`
	}
}
