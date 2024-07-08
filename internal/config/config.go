package config

type Server struct {
	Address string `yaml:"address"`
}

type App struct {
	DbUrl string `yaml:"db_url"`
}

type Config struct {
	Server `yaml:"server"`
	App    `yaml:"app"`
}
