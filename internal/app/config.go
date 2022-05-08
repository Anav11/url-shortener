package app

type Config struct {
	ServerAddress	string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL			string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath	string `env:"FILE_STORAGE_PATH" envDefault:"./urls_db.csv"`
	SecretKey		string `env:"SECRET_KEY" envDefault:"pf-ctvm-gtxfnzvb"`
	DatabaseDSN		string `env:"DATABASE_DSN" envDefault:"postgres://a.averin:password@localhost:7070/shortener"`
}
