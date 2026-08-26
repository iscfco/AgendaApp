package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de GORM que usarán tus repositorios [cite: 67, 72]
var DB *gorm.DB

func GetDB() *gorm.DB {
	if DB == nil {
		panic("No se ha inicializado la base de datos")
	}
	return DB
}

// InitDB inicializa la conexión a MySQL usando GORM
func InitDB(user, pass, host, port, dbName string) *gorm.DB {
	// Formato DSN para MySQL [cite: 42]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, pass, host, port, dbName)

	var err error
	// Inicializar la conexión con configuración de GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Útil para depurar queries [cite: 48]
	})

	if err != nil {
		log.Fatalf("Error al conectar a la base de datos con GORM: %v", err)
	}

	// Configuración del Pool de conexiones subyacente (sql.DB)
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener la instancia sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)                 // Máximo de conexiones abiertas
	sqlDB.SetMaxIdleConns(25)                 // Máximo de conexiones inactivas
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Tiempo de vida de conexión [cite: 42]

	fmt.Println("Conexión exitosa a MySQL mediante GORM")

	return DB
}
