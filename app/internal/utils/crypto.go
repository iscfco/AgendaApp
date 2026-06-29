package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword recibe una contraseña en texto plano y devuelve su hash
func HashPassword(password string) (string, error) {
	// GenerateFromPassword crea un hash con un "costo" (por defecto 10)
	// El costo determina qué tan lento es el proceso (más lento = más seguro)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara una contraseña en texto plano con el hash guardado en la BD
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
