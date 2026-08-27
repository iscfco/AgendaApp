package middlewares

import (
	"agenda-app/app/internal/utils"
	"agenda-app/app/internal/utils/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica si existe la cookie de sesión
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(utils.SessionCookieHeader)

		// Si hay un error (la cookie no existe o expiró)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// 2. Verificamos si ese ID realmente existe en nuestro mapa del servidor
		_, existe := sessions.ObtenerSessionDataPorSesion(sessionID)
		if !existe {
			// El ID es falso o el servidor se reinició y la memoria se borró
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
