package controllers

import (
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/utils"
	"agenda-app/app/internal/utils/logs"
	"agenda-app/app/internal/utils/sessions"
	"agenda-app/app/internal/views"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HomeController struct {
	//service *services.LoginService
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

// Create maneja la petición GET /
func (ctrl *HomeController) GetHome(c *gin.Context) {
	ctx := c.Request.Context()
	logs.Logger(ctx).Info("getting home")

	// Unimos el html base con el html de la vista
	tmpl, err := template.New("base").ParseFS( // tempate.New("base") debe se igual al nombre del html base
		views.ViewsFS,
		"layout/base.html", // Archivo base con {{template "control-buttons" .}} y {{template "content" .}}
		"orders/list.html", // El archivo que define {{define "control-buttons"}} y {{define "content"}}
	)
	if err != nil {
		logs.Logger(ctx).Error("Error cargando plantillas", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error cargando plantillas")
		return
	}

	// Obtenemos los datos del usuario desde la sesion
	tokenDeSesion, err := c.Cookie(utils.SessionCookieHeader)
	if err != nil { // El error significa que la cookie no existe o expiró
		logs.Logger(ctx).Error("No se pudo obtener el token de sesión", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener el token de sesión")
		return
	}

	data, ok := sessions.ObtenerSessionDataPorSesion(tokenDeSesion)
	if !ok {
		logs.Logger(ctx).Error("No se pudo obtener los datos de sesión de la cache de sesiones")
		c.String(http.StatusInternalServerError, "No se pudo obtener los datos de sesión")
		return
	}
	user, ok := data.(models.User)
	if !ok {
		logs.Logger(ctx).Error("Error al hacer casting a User", zap.Any("data", user))
		c.String(http.StatusInternalServerError, "No se pudo obtener los datos por error interno")
		return
	}

	// Ejecutamos directamente con .Execute porque el template ya sabe que su nodo principal es "base"
	err = tmpl.Execute(c.Writer, gin.H{
		"title":    "Agenda App - Pedidos",
		"fullname": user.UserFullName,
		"email":    user.Email,
		"role":     user.Role,
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}
