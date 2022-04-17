package app

type Config struct {
	ServerAddress	string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL		string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath	string `enb:"FILE_STORAGE_PATH" envDefault:"./urls_db.csv"`
}