package config

type Config struct {
	MongoURI string
	Port     string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI: "mongodb://root:root@mongo:27017/app?authSource=admin",
	}
}
