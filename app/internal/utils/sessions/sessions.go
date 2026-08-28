package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var (
	sessionsMap = make(map[string]interface{})
	mutex       sync.RWMutex
)

// GenerarSessionID crea un token aleatorio, único e imposible de adivinar
func GenerarSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GuardarSesion anota el ID de sesión asociado a la data que queremos guardar
func GuardarSesion(sessionID string, sessionData interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	sessionsMap[sessionID] = sessionData
}

// ObtenerSessionDataPorSesion busca si el ID existe en el mapa y devuelve la data
func ObtenerSessionDataPorSesion(sessionID string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	sessionData, existe := sessionsMap[sessionID]
	return sessionData, existe
}

// EliminarSesion borra el registro (útil para el Logout)
func EliminarSesion(sessionID string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(sessionsMap, sessionID)
}
