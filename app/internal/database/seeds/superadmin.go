package database

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// EnsureSuperAdmin verifica y crea el usuario administrador inicial
func EnsureSuperAdmin(db *sql.DB) {
	// 1. Obtener credenciales desde variables de entorno (.env)
	email := os.Getenv("SUPERADMIN_EMAIL") // TODO: enhance config management
	password := os.Getenv("SUPERADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("⚠️ Advertencia: Variables de entorno de SuperAdmin no configuradas. Saltando seeding.")
		return
	}

	// 2. Verificar si el usuario ya existe en la base de datos
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Fatalf("❌ Error al verificar existencia de SuperAdmin: %v", err)
	}

	if exists {
		log.Println("✅ El SuperAdmin ya existe en la base de datos.")
		return
	}

	// 3. Si no existe, Hashear la contraseña usando tu función de Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Error al hashear la contraseña del SuperAdmin: %v", err)
	}

	// 4. Insertar el SuperAdmin en la base de datos
	insertQuery := `
		INSERT INTO users (email, password, role) 
		VALUES ($1, $2, 'superadmin')`

	_, err = db.Exec(insertQuery, email, string(hashedPassword))
	if err != nil {
		log.Fatalf("❌ Error al insertar el SuperAdmin: %v", err)
	}

	log.Println("🚀 SuperAdmin creado exitosamente por primera vez.")
}
