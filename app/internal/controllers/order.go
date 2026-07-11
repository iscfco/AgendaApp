package controllers

import (
	"agenda-app/app/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController(s *services.OrderService) *OrderController {
	return &OrderController{service: s}
}

// Create maneja la petición POST /orders
func (ctrl *OrderController) Create(c *gin.Context) {
	// var newOrder models.Order

	// // 1. "Bind" del JSON: Valida que el body traiga los campos del struct
	// if err := c.ShouldBindJSON(&newOrder); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del pedido inválidos: " + err.Error()})
	// 	return
	// }

	// // 2. Obtener el 'Actor' (el usuario que hace la petición)
	// // Generalmente extraído de un Middleware de Auth previo
	// actor, exists := c.Get("currentUser")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no identificado"})
	// 	return
	// }

	// // 3. Llamar al servicio para aplicar la lógica de negocio
	// // Pasamos el actor (usuario A) y los datos del nuevo pedido (usuario B/cliente)
	// err := ctrl.service.CreateOrder(actor.(models.User), newOrder)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// // 4. La Vista: Respuesta de éxito
	// c.JSON(http.StatusCreated, gin.H{
	// 	"message": "Pedido creado exitosamente",
	// 	"data":    newOrder,
	// })
}
