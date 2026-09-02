package controllers

import (
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/utils"
	"agenda-app/app/internal/utils/sessions"
	"fmt"

	"github.com/gin-gonic/gin"
)

func getUserFromSession(c *gin.Context) (user models.User, err error) {
	var tokenDeSesion string

	// Obtenemos los datos del usuario desde la sesion
	tokenDeSesion, err = c.Cookie(utils.SessionCookieHeader)
	if err != nil { // El error significa que la cookie no existe o expiró
		return user, fmt.Errorf("No se pudo obtener el token de sesión: %w", err)
	}

	data, ok := sessions.ObtenerSessionDataPorSesion(tokenDeSesion)
	if !ok {
		return user, fmt.Errorf("No se pudo obtener los datos de sesión de la cache de sesiones: %w", err)
	}

	user, ok = data.(models.User)
	if !ok {
		return user, fmt.Errorf("No se pudo obtener los datos por error interno: %w", err)
	}

	return user, nil
}
