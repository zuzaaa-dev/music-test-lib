package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env        string `env:"ENV" env-default:"development"`
	HTTPServer HTTPServer
	DataBase   DataBase
	API        API
}

type DataBase struct {
	Host           string `env:"DB_HOST" env-required:"true"`
	Port           string `env:"DB_PORT" env-required:"true"`
	User           string `env:"DB_USER" env-required:"true"`
	Password       string `env:"DB_PASSWORD" env-required:"true"`
	Name           string `env:"DB_NAME" env-required:"true"`
	FileMigrations string `env:"DB_FILE_MIGRATIONS" env-required:"true"`
}

type HTTPServer struct {
	Address string `env:"HTTP_SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
}

type API struct {
	MusicInfoURL string `env:"API_MUSIC_INFO_URL" env-required:"true"`
}

func MustLoad() *Config {
	var config Config
	// Загружаем переменные окружения из .env
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("Ошибка загрузки файла конфигураций: %s", err)
	}
	return &config
}

/*В случае если файл с конфигурациями в yaml*/
//func fetchConfigPath(filename string) string {
//	path, err := filepath.Abs(filename)
//	if err != nil {
//		errString := fmt.Sprintf("Error getting absolute path of config file %s: %s\n", filename, err)
//		panic(errString)
//	}
//	return path
//}
//func MustLoad(filename string) *Config {
//	// Установка пути до файла конфигурации
//	configPath := fetchConfigPath(filename)
//	if configPath == "" {
//		panic("config path not set")
//	}
//	fmt.Println(configPath)
//	if q, err := os.Stat(configPath); os.IsNotExist(err) {
//		fmt.Println(q)
//		panic("config path does not exist")
//	}
//	config := &Config{}
//	if err := cleanenv.ReadConfig(configPath, config); err != nil {
//		panic(err)
//	}
//	return config
//}
