package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"agenda-app/app/internal/controllers"
	"agenda-app/app/internal/middlewares"
	"agenda-app/app/internal/repository"
	"agenda-app/app/internal/services"
	"agenda-app/app/internal/views" // Ajusta según tu go.mod

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Iniciando el servidor...")

	// Cargar configs // TODO: encriptarlos

	// DB connection

	// Call DB Seed

	// Dependencias
	userRepo := repository.NewUserRepository(nil)

	loginSvc := services.NewLoginService(userRepo)

	home := controllers.NewHomeController()
	login := controllers.NewLoginController(loginSvc)

	// Crear router
	r := gin.Default()

	// Static Files
	staticSubFS, err := fs.Sub(views.ViewsFS, "static")
	if err != nil {
		panic("No se pudo encontrar la carpeta static dentro del embed: " + err.Error())
	}
	r.StaticFS("/static", http.FS(staticSubFS))

	// Public router
	loginRoutes := r.Group("/login")
	{
		loginRoutes.GET("/", login.GetLogin)
		loginRoutes.POST("/", login.Login)
	}

	// Private router
	r.Use(middlewares.AuthMiddleware())

	r.GET("/", home.GetHome)

	r.GET("/orders/new", func(c *gin.Context) {
		tmpl, err := template.New("base").ParseFS(
			views.ViewsFS,
			"layout/base.html",
			"orders/create.html",
		)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error cargando plantillas: %v", err)
			return
		}

		// Ejecutamos el template apuntando a "base"
		err = tmpl.Execute(c.Writer, gin.H{
			"title": "Registro de Órdenes",
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error renderizando: %v", err)
		}
	})

	r.Run(":8080")
}
