package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"agenda-app/app/config"
	"agenda-app/app/internal/controllers"
	"agenda-app/app/internal/database"
	"agenda-app/app/internal/database/bootstrap"
	"agenda-app/app/internal/middlewares"
	"agenda-app/app/internal/repository"
	"agenda-app/app/internal/services"
	"agenda-app/app/internal/utils/logs"
	"agenda-app/app/internal/views"
	"agenda-app/app/migrations"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Iniciando el servidor...")

	// Setup logger
	logger := logs.InitLogger()

	// Cargar configs // TODO: encriptarlos
	cfg := config.LoadConfig()

	// DB connection
	db, err := database.InitDB(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)
	if err != nil {
		logger.Fatal("No se pudo conectar a la base de datos", zap.String("err", err.Error()))
	}
	logger.Info("Conectado a la BD")

	// Run migrations
	if err := migrations.Run(db); err != nil {
		logger.Fatal("No se pudieron correr las migraciones", zap.String("err", err.Error()))
	}

	// Bootstrap admin
	err = bootstrap.Run(db, cfg.Bootstrap.SuperAdmin.Email, cfg.Bootstrap.SuperAdmin.Password, cfg.Bootstrap.SuperAdmin.Name)
	if err != nil {
		logger.Fatal("No se pudo crear el super admin", zap.String("err", err.Error()))
	}

	// Dependencias
	// - Repos
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// - Services
	userSvc := services.NewUserService(userRepo)
	orderSvc := services.NewOrderService(orderRepo, userRepo)

	// - Controllers
	order := controllers.NewOrderController(orderSvc, userSvc)
	login := controllers.NewLoginController(userSvc)

	// Run API
	{
		// Crear router
		r := gin.Default()

		r.Use(logs.GinZapMiddleware(logger))

		// Static Files
		staticSubFS, err := fs.Sub(views.ViewsFS, "static")
		if err != nil {
			panic("No se pudo encontrar la carpeta static dentro del embed: " + err.Error())
		}
		r.StaticFS("/static", http.FS(staticSubFS))

		// Public routes
		loginRoutes := r.Group("/login")
		{
			loginRoutes.GET("/", login.GetLogin)
			loginRoutes.POST("/", login.Login)
		}

		// Private router
		r.Use(middlewares.AuthMiddleware())
		r.GET("/", order.GetOrderView)
		r.GET("/order/view", order.GetCreateOrderView)
		r.POST("/order", order.Create)
		r.GET("/order/:id/view", order.GetOrderDetailsView)
		r.PUT("/order/:id", order.Update)
		r.DELETE("/order/:id", order.Delete)

		r.Run(":8080")
	}
}
