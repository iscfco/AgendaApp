package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var (
	sessionsMap = make(map[string]string)
	mutex       sync.RWMutex
)

// GenerarSessionID crea un token aleatorio, único e imposible de adivinar
func GenerarSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GuardarSesion anota el ID de sesión asociado al email en nuestro mapa
func GuardarSesion(sessionID string, email string) {
	mutex.Lock()
	defer mutex.Unlock()

	sessionsMap[sessionID] = email
}

// ObtenerUsuarioPorSesion busca si el ID existe en el mapa y devuelve el email
func ObtenerUsuarioPorSesion(sessionID string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	email, existe := sessionsMap[sessionID]
	return email, existe
}

// EliminarSesion borra el registro (útil para el Logout)
func EliminarSesion(sessionID string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(sessionsMap, sessionID)
}
