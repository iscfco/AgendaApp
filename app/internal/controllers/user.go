package controllers

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/models/filters"
	"agenda-app/app/internal/services"
	"agenda-app/app/internal/utils/logs"
	"agenda-app/app/internal/views"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userSvc services.UserService
}

func NewUserController(userSvc services.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

func checkRole(requester models.User) error {
	if requester.Role != models.UserRoleAdmin && requester.Role != models.UserRoleSuperAdmin {
		return errorhandling.ErrForbidden
	}

	return nil
}

// Create maneja la petición GET /user
func (ctrl *UserController) GetUsersView(c *gin.Context) {
	ctx := c.Request.Context()

	requester, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("Error al obtener el requester", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al obtener el requester")
		return
	}

	if err := checkRole(requester); err != nil {
		c.String(http.StatusForbidden, "Permisos insuficientes")
		return
	}

	// Unimos el html base con el html de la vista
	tmpl, err := template.New("base").ParseFS(
		views.ViewsFS,
		"layout/base.html",
		"users/list-user.html",
	)
	if err != nil {
		logs.Logger(ctx).Error("Error cargando plantillas", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error cargando plantillas")
		return
	}

	// Get Users
	var filter filters.GetUsers
	if err := c.ShouldBindQuery(&filter); err != nil {
		logs.Logger(ctx).Error("No se pudo obtener los filtros", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener los filtros")
		return
	}

	var users []models.User
	var totalRecords int64
	var totalPages int
	var pages []int
	if filter.Show {
		// Buscar desde la base de datos
		users, totalRecords, err = ctrl.userSvc.ListUsers(requester, filter)
		if err != nil {
			logs.Logger(ctx).Error("Error al obtener los pedidos", zap.Error(err))
			c.String(http.StatusInternalServerError, "Error al obtener los pedidos")
			return
		}
		pages, totalPages = generatePages(int(totalRecords), filter.Limit, filter.Page)
	}

	err = tmpl.Execute(c.Writer, gin.H{
		"title":      "Agenda App - Pedidos",
		"userActive": true,

		// Datos del usuario
		"fullname": requester.UserFullName,
		"email":    requester.Email,
		"role":     requester.Role,

		// Datos de las ordenes
		"show":            filter.Show,
		"users":           users,
		"orderStatusText": models.OrderStatusText,

		// Paginacion
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"pages":         pages,
		"current_page":  filter.Page,
		"limit":         filter.Limit,

		// Returning back all filters
		"filter_username": filter.UserFullName,
		"filter_role":     filter.Role,
		"filter_email":    filter.Email,
		"filter_status":   filter.Status,
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}

// Create maneja la petición POST /user
func (ctrl *UserController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	requester, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("Error al obtener el actor", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := checkRole(requester); err != nil {
		c.String(http.StatusForbidden, "Permisos insuficientes")
		return
	}

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		logs.Logger(ctx).Error("Datos del usuario inválidos", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del usuario inválidos: " + err.Error()})
		return
	}

	_, err = ctrl.userSvc.RegisterNewUser(requester, newUser)
	if err != nil {
		logs.Logger(ctx).Error("Error al crear el usuario", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// // 4. La Vista: Respuesta de éxito
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"data":    nil,
	})
}

func (ctrl *UserController) GetCreateUserView(c *gin.Context) {
	ctx := c.Request.Context()

	requester, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener el usuario", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener la sesion del usuario")
		return
	}

	err = checkRole(requester)
	if err != nil {
		c.String(http.StatusForbidden, "Permisos insuficientes")
		return
	}

	tmpl, err := template.New("base").ParseFS(
		views.ViewsFS,
		"layout/base.html",
		"users/create-user.html",
	)
	if err != nil {
		logs.Logger(ctx).Error("Error cargando plantillas", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error cargando plantillas")
		return
	}

	err = tmpl.Execute(c.Writer, gin.H{
		"title": "Agenda App - Crear Pedidos",
		// Datos del usuario
		"fullname": requester.UserFullName,
		"email":    requester.Email,
		"role":     requester.Role,
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}

func (ctrl *UserController) GetUserDetailsView(c *gin.Context) {
	ctx := c.Request.Context()

	// Get requester
	requester, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener el usuario", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener la sesion del usuario")
		return
	}

	if err := checkRole(requester); err != nil {
		c.String(http.StatusForbidden, "Permisos insuficientes")
		return
	}

	// Unimos el html base con el html de la vista
	tmpl, err := template.New("base").ParseFS(
		views.ViewsFS,
		"layout/base.html",
		"users/user-details.html",
	)
	if err != nil {
		logs.Logger(ctx).Error("Error cargando plantillas", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error cargando plantillas")
		return
	}

	// Get order
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener el id de la orden", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener el id de la orden")
		return
	}

	user, err := ctrl.userSvc.GetUserById(uint(idInt))
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener la orden", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener la orden")
		return
	}

	var showUpdateButtons bool
	if requester.Role == "superadmin" {
		showUpdateButtons = true
	}

	if requester.Role == "admin" && user.Role != "superadmin" {
		showUpdateButtons = true
	}

	err = tmpl.Execute(c.Writer, gin.H{
		"title": "Agenda App - Crear Pedidos",
		// Datos del usuario
		"fullname": requester.UserFullName,
		"email":    requester.Email,
		"role":     requester.Role,

		// Datos de la orden
		"user_id":         idInt,
		"userfullname":    user.UserFullName,
		"user_email":      user.Email,
		"user_phone":      user.Phone,
		"user_role":       user.Role,
		"user_status":     user.Status,
		"user_created_at": user.CreatedAt.Format(time.DateTime),

		"show_update_buttons": showUpdateButtons,
	})
	if err != nil {
		logs.Logger(ctx).Error("Error renderizando", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error renderizando")
	}
}

// Create maneja la petición PUT /user/:id
func (ctrl *UserController) Update(c *gin.Context) {
	ctx := c.Request.Context()

	// Get requester
	requester, err := getUserFromSession(c)
	if err != nil {
		logs.Logger(ctx).Error("No se pudo obtener el usuario", zap.Error(err))
		c.String(http.StatusInternalServerError, "No se pudo obtener la sesion del usuario")
		return
	}

	if err := checkRole(requester); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permisos insuficientes"})
		return
	}

	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		logs.Logger(ctx).Error("ID del usuario inválido", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del usuario inválido: " + err.Error()})
		return
	}

	var updates models.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		logs.Logger(ctx).Error("Datos del usuario inválidos", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del usuario inválidos: " + err.Error()})
		return
	}

	updates.ID = uint(userIDInt)
	err = ctrl.userSvc.UpdateUser(requester, updates)
	if err != nil {
		logs.Logger(ctx).Error("Error al actualizar el usuario", zap.Error(err))
		errMsg := "no se puedo actualizar el usuario"
		switch {
		case errors.Is(err, errorhandling.ErrRole):
			errMsg = "No tienes permiso para actualizar el rol del usuario"
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// // 4. La Vista: Respuesta de éxito
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pedido actualizado exitosamente",
	})
}
