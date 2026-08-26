package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"agenda-app/app/config"
	"agenda-app/app/internal/controllers"
	"agenda-app/app/internal/database"
	"agenda-app/app/internal/middlewares"
	"agenda-app/app/internal/repository"
	"agenda-app/app/internal/services"
	"agenda-app/app/internal/utils/logs"
	"agenda-app/app/internal/views" // Ajusta según tu go.mod

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Iniciando el servidor...")

	// Cargar configs // TODO: encriptarlos
	cfg := config.LoadConfig()

	// DB connection
	db := database.InitDB(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)

	// Call DB Seed

	// Dependencias
	// - Repos
	userRepo := repository.NewUserRepository(db)

	// - Services
	userSvc := services.NewUserService(userRepo)

	// - Controllers
	home := controllers.NewHomeController()
	login := controllers.NewLoginController(userSvc)

	// Crear router
	r := gin.Default()

	// Setup logger
	logger := logs.InitLogger()
	r.Use(logs.GinZapMiddleware(logger))

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
