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
		email, existe := sessions.ObtenerUsuarioPorSesion(sessionID)
		if !existe {
			// El ID es falso o el servidor se reinició y la memoria se borró
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// 3. Opcional: Guardamos el email del usuario en el contexto de Gin
		// por si tus otras rutas (como /dashboard) quieren saber qué usuario inició sesión.
		c.Set("user_email", email)

		c.Next()
	}
}
