// internal/services/order_service.go
package services

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type OrderService interface {
	CreateOrder(requestor models.User, order models.Order)
	ListOrders(query models.Order) ([]models.Order, error)
	UpdateOrder(requestor models.User, model models.Order) error
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repo: r}
}

type orderService struct {
	repo repository.OrderRepository
}

func (s *orderService) CreateOrder(requestor models.User, order models.Order) {
	order.AuthorID = requestor.ID
	order.CreatedAt = requestor.CreatedAt
	order.Status = models.OrderStatusPending

	s.repo.Create(order)
}

func (s *orderService) ListOrders(query models.Order) ([]models.Order, error) {
	query = models.Order{
		AuthorID:     query.AuthorID,
		ClientName:   query.ClientName,
		CreatedAt:    query.CreatedAt,
		DeliveryDate: query.DeliveryDate,
		Status:       query.Status,
	}

	result, err := s.repo.ReadByQuery(query)
	if err != nil {
		return nil, fmt.Errorf("%w: error al listar pedidos: %v", errorhandling.ErrInternal)
	}

	return result, nil
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
		updates.ChangeLog = string(newChangeLog)
	}

	// Actualizar
	return s.repo.Update(updates)
}
