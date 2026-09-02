package controllers

import (
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/models/filters"
	"agenda-app/app/internal/services"
	"agenda-app/app/internal/utils"
	"agenda-app/app/internal/utils/logs"
	"agenda-app/app/internal/utils/sessions"
	"agenda-app/app/internal/views"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderController struct {
	service services.OrderService
}

func NewOrderController(s services.OrderService) *OrderController {
	return &OrderController{service: s}
}

// Create maneja la petición POST /orders
func (ctrl *OrderController) Create(c *gin.Context) {
	var newOrder models.Order
	ctx := c.Request.Context()

	// 1. "Bind" del JSON: Valida que el body traiga los campos del struct
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		logs.Logger(ctx).Error("Datos del pedido inválidos", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del pedido inválidos: " + err.Error()})
		return
	}

	// 2. Obtener el 'autor' (el usuario que hace la petición)
	user, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("Error al obtener el actor", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Llamar al servicio para aplicar la lógica de negocio
	// Pasamos el actor (usuario A) y los datos del nuevo pedido (usuario B/cliente)
	err = ctrl.service.CreateOrder(user, newOrder)
	if err != nil {
		logs.Logger(ctx).Error("Error al crear el pedido", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// // 4. La Vista: Respuesta de éxito
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pedido creado exitosamente",
		"data":    nil,
	})
}

// Create maneja la petición GET /
func (ctrl *OrderController) GetOrderView(c *gin.Context) {
	ctx := c.Request.Context()
	logs.Logger(ctx).Info("getting orders")

	// Unimos el html base con el html de la vista
	tmpl, err := template.New("base").ParseFS( // tempate.New("base") debe se igual al nombre del html base
		views.ViewsFS,
		"layout/base.html",        // Archivo base con {{template "control-buttons" .}} y {{template "content" .}}
		"orders/list-orders.html", // El archivo que define {{define "control-buttons"}} y {{define "content"}}
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

	// Get Orders
	var filter filters.GetOrders
	if err := c.ShouldBindQuery(&filter); err != nil {
		logs.Logger(ctx).Error("No se pudo obtener los filtros", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener los filtros")
		return
	}

	logs.Logger(ctx).Info("Payload received", zap.Any("payload", filter))
	var orders []models.Order
	var totalRecords int64
	var totalPages int
	var pages []int
	if filter.Show {
		// Buscar desde la base de datos
		orders, totalRecords, err = ctrl.service.ListOrders(ctx, filter)
		if err != nil {
			logs.Logger(ctx).Error("Error al obtener los pedidos", zap.Error(err))
			c.String(http.StatusInternalServerError, "Error al obtener los pedidos")
			return
		}
		if len(orders) > 0 {
			totalPages = int(totalRecords) / filter.Limit
			if int(totalRecords)%filter.Limit != 0 || totalPages == 0 {
				totalPages++ // Si sobran registros, añade una página extra
			}
		}

		// Generate pages
		if totalPages <= 7 {
			// Generate from 1 to totalPages
			pages = make([]int, totalPages)
			for i := range totalPages {
				pages[i] = i + 1
			}
		} else {
			pages = make([]int, 7)
			// Calculate starting and ending page
			startingPage := filter.Page - 3
			endingPage := filter.Page + 3

			if endingPage > totalPages { // Verify if we can add more pages >>
				endingPage = totalPages
				startingPage = totalPages - 6
			}

			if startingPage < 1 { // Verify if we can add more pages <<
				startingPage = 1
				endingPage = 7
			}

			// Generate from startingPage to endingPage
			for i := 0; i < 7; i++ {
				pages[i] = startingPage + i
			}
		}
	}

	// Ejecutamos directamente con .Execute porque el template ya sabe que su nodo principal es "base"
	err = tmpl.Execute(c.Writer, gin.H{
		"title": "Agenda App - Pedidos",
		// Datos del usuario
		"fullname": user.UserFullName,
		"email":    user.Email,
		"role":     user.Role,

		// Datos de las ordenes
		"show":            filter.Show,
		"orders":          orders,
		"orderStatusText": models.OrderStatusText,

		// Paginacion
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"pages":         pages,
		"current_page":  filter.Page,
		"limit":         filter.Limit,

		// Returning back all filters
		"filter_user_creator_name": filter.UserCreatorName,
		"filter_keyword":           filter.Keyword,
		"filter_client_name":       filter.ClientName,
		"filter_from_date":         filter.From,
		"filter_to_date":           filter.To,
		"filter_status":            models.OrderStatusText[models.OrderStatus(filter.Status)],
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}

func (ctrl *OrderController) GetCreateOrderView(c *gin.Context) {
	ctx := c.Request.Context()

	// Unimos el html base con el html de la vista
	tmpl, err := template.New("base").ParseFS( // tempate.New("base") debe se igual al nombre del html base
		views.ViewsFS,
		"layout/base.html",         // Archivo base con {{template "control-buttons" .}} y {{template "content" .}}
		"orders/create-order.html", // El archivo que define {{define "control-buttons"}} y {{define "content"}}
	)
	if err != nil {
		logs.Logger(ctx).Error("Error cargando plantillas", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error cargando plantillas")
		return
	}

	// Get User
	user, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener el usuario", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener la sesion del usuario")
		return
	}

	// Ejecutamos directamente con .Execute porque el template ya sabe que su nodo principal es "base"
	err = tmpl.Execute(c.Writer, gin.H{
		"title": "Agenda App - Crear Pedidos",
		// Datos del usuario
		"fullname": user.UserFullName,
		"email":    user.Email,
		"role":     user.Role,
		// Datos inyectados al form de crear orden
		"author": user.UserFullName,
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}
