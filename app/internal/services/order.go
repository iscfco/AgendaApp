// internal/services/order_service.go
package services

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/models/filters"
	"agenda-app/app/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type OrderService interface {
	CreateOrder(requestor models.User, order models.Order) error
	ListOrders(ctx context.Context, filter filters.GetOrders) ([]models.Order, int64, error)
	UpdateOrder(requestor models.User, model models.Order) error
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repo: r}
}

type orderService struct {
	repo repository.OrderRepository
}

func (s *orderService) CreateOrder(requestor models.User, order models.Order) error {
	return s.repo.Create(models.Order{
		AuthorID:      requestor.ID,
		ClientName:    order.ClientName,
		ClientPhone:   order.ClientPhone,
		ClientAddress: order.ClientAddress,
		TotalPrice:    order.TotalPrice,
		DownPayment:   order.DownPayment,
		// CreatedAt:           requestor.CreatedAt, Automatico
		// UpdatedAt:           requestor.CreatedAt,
		UpdatedBy:    requestor.ID,
		DeliveryDate: order.DeliveryDate,
		Description:  order.Description,
		Status:       models.OrderStatusPending,
		//ChangeLog:           datatypes.JSON{},
		//StoredInChangeLogAt: time.Time{},
	})
}

func (s *orderService) ListOrders(ctx context.Context, filter filters.GetOrders) ([]models.Order, int64, error) {
	result, total, err := s.repo.ReadByQuery(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: error al listar pedidos: %v", errorhandling.ErrInternal, err)
	}

	return result, total, nil
}

func (s *orderService) UpdateOrder(requestor models.User, model models.Order) error {
	// Validar que el pedido exista
	orderInDB, err := s.repo.ReadByID(model.ID)
	if errors.Is(err, errorhandling.ErrNotFoundError) {
		return fmt.Errorf("%w: el pedido con id '%d' no existe", err, model.ID)
	}
	if err != nil {
		return fmt.Errorf("%w: error intentando encontrar pedido a modificar: %v", errorhandling.ErrInternal, err)
	}

	// Preparar updates
	updates := models.Order{
		ClientName:    model.ClientName,
		ClientPhone:   model.ClientPhone,
		ClientAddress: model.ClientAddress,
		TotalPrice:    model.TotalPrice,
		DownPayment:   model.DownPayment,
		Description:   model.Description,
		DeliveryDate:  model.DeliveryDate,
		Status:        model.Status,
		UpdatedAt:     time.Now(),
		UpdatedBy:     requestor.ID,
	}

	// Definir delivery date
	{
		if model.Status == models.OrderStatusDelivered {
			updates.DeliveryDate = time.Now()
		}
		if model.Status == models.OrderStatusPending {
			updates.DeliveryDate = time.Time{}
		}
	}

	// - Preparar change log
	{
		// Convertir string a arreglo de orders
		currentChangeLog := []models.Order{}
		err = json.Unmarshal([]byte(orderInDB.ChangeLog), &currentChangeLog)
		if err != nil {
			return fmt.Errorf("%w: error al deserializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Agregar order al arreglo
		orderInDB.StoredInChangeLogAt = time.Now()
		currentChangeLog = append(currentChangeLog, orderInDB)

		// Convertir de nuevo a string
		newChangeLog, err := json.Marshal(currentChangeLog)
		if err != nil {
			return fmt.Errorf("%w: error al serializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Actualizar change log
		updates.ChangeLog = datatypes.JSON(newChangeLog)
	}

	// Actualizar
	return s.repo.Update(updates)
}
