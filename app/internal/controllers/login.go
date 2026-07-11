package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"agenda-app/app/internal/services"
	"agenda-app/app/internal/utils/sessions"
	"agenda-app/app/internal/views"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	service services.LoginService
}

func NewLoginController(s services.LoginService) *LoginController {
	return &LoginController{service: s}
}

// Create maneja la petición GET /login
func (ctrl *LoginController) GetLogin(c *gin.Context) {
	// tmpl, err := template.New("base").ParseFS(
	// 	views.ViewsFS,
	// 	"layout/base.html",
	// 	"orders/create.html",
	// )
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, "Error cargando plantillas: %v", err)
	// 	return
	// }

	// // Ejecutamos el template apuntando a "base"
	// err = tmpl.Execute(c.Writer, gin.H{
	// 	"title": "Registro de Órdenes",
	// })
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, "Error renderizando: %v", err)
	// }

	// Leemos el archivo del login desde tu embed
	htmlBytes, err := views.ViewsFS.ReadFile("login/login.html") // Ajusta la ruta a tu archivo real
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al cargar la página")
		return
	}

	// Parseamos la plantilla al vuelo
	tmpl, err := template.New("login").Parse(string(htmlBytes))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al procesar la plantilla")
		return
	}

	// Le indicamos al navegador que es un HTML
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")

	// Ejecutamos pasando valores VACÍOS (gin.H{}) para que el {{if .error}} se evalúe como falso
	// y desaparezca por completo de la pantalla.
	tmpl.Execute(c.Writer, gin.H{})
}

func (ctrl *LoginController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// 1. Crear una función interna para renderizar el HTML del embed con datos
	renderLoginWithError := func(status int, mensajeError string) {
		// Ajusta la ruta exacta según dónde guardaste tu archivo (ej: "login/login.html" o "index.html")
		htmlBytes, err := views.ViewsFS.ReadFile("login/login.html")
		if err != nil {
			fmt.Println("Error leyendo embed:", err)
			c.String(http.StatusInternalServerError, "Error al cargar la página de login")
			return
		}

		// Creamos una plantilla al vuelo con el contenido del archivo
		tmpl, err := template.New("login").Parse(string(htmlBytes))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error parseando plantilla")
			return
		}

		// Enviamos el HTML ejecutado directamente al navegador
		c.Status(status)
		c.Header("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(c.Writer, gin.H{
			"error":    mensajeError,
			"oldEmail": email, // Opcional: para no borrarle el correo que ya escribió
		})
	}

	// 2. Validación de campos vacíos
	if email == "" || password == "" {
		renderLoginWithError(http.StatusBadRequest, "Todos los campos son obligatorios")
		return
	}

	// 3. Validación de credenciales
	// TODO: Move to DB queries
	if !(email == "admin@correo.com" && password == "123456") {
		// Login fallido: renderizamos usando nuestra función utilitaria
		renderLoginWithError(http.StatusUnauthorized, "Credenciales incorrectas")
		return
	}

	// 4. Guardar la sesión
	tokenDeSesion := sessions.GenerarSessionID()
	sessions.GuardarSesion(tokenDeSesion, email)

	// Guardamos la cookie en el navegador del usuario
	c.SetCookie(
		"session_token", // Nombre de la cookie (debe coincidir con tu Middleware)
		tokenDeSesion,   // Valor que se guardará
		60*60*24,        // Tiempo de vida en SEGUNDOS ( 24 hrs = 60 segundos * 60 minutos * 24 horas )
		"/",             // Ruta donde la cookie es válida ("/" significa en toda la app)
		"localhost",     // Dominio de tu servidor
		false,           // Secure: true solo para HTTPS.
		true,            // HttpOnly: Evita que JavaScript robe la cookie (Protección XSS).
	)

	c.Redirect(http.StatusSeeOther, "/dashboard")
}
