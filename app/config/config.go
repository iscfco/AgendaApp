package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config almacena las variables de entorno de la aplicación
type Config struct {
	Port  string
	DBURI string
	Env   string
	DB    DB
}

type DB struct {
	User   string
	Pass   string
	Host   string
	Port   string
	DBName string
}

// LoadConfig lee el archivo .env y mapea los valores
func LoadConfig() *Config {
	// Intenta cargar el archivo .env (si no existe, continuará buscando en el sistema)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró el archivo .env, usando variables de entorno del sistema")
	}

	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBURI: getEnv("DATABASE_URL", ""),
		Env:   getEnv("APP_ENV", "development"),
		DB: DB{
			User:   getEnv("DB_USER", ""),
			Pass:   getEnv("DB_PASS", ""),
			Host:   getEnv("DB_HOST", ""),
			Port:   getEnv("DB_PORT", ""),
			DBName: getEnv("DB_NAME", ""),
		},
	}
}

// getEnv es una función de ayuda para asignar valores por defecto si la variable está vacía
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
