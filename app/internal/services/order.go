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
	GetOrderById(id uint) (models.Order, error)
	DeleteOrder(id uint) error
}

func NewOrderService(r repository.OrderRepository, userRepo repository.UserRepository) OrderService {
	return &orderService{repo: r, userRepo: userRepo}
}

type orderService struct {
	repo     repository.OrderRepository
	userRepo repository.UserRepository
}

func (s *orderService) GetOrderById(id uint) (models.Order, error) {
	return s.repo.ReadByID(id)
}

func (s *orderService) CreateOrder(requestor models.User, order models.Order) error {
	return s.repo.Create(models.Order{
		AuthorID:      requestor.ID,
		ClientName:    order.ClientName,
		ClientPhone:   order.ClientPhone,
		ClientAddress: order.ClientAddress,
		TotalPrice:    order.TotalPrice,
		DownPayment:   order.DownPayment,
		UpdatedBy:     requestor.ID,
		DeliveryDate:  order.DeliveryDate,
		Description:   order.Description,
		Status:        models.OrderStatusPending,
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
		ID:            model.ID,
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

	// - Preparar change log
	{
		if orderInDB.ChangeLog == nil {
			orderInDB.ChangeLog = datatypes.JSON("[]")
		}

		// Convertir string a arreglo de orders
		currentChangeLog := []map[string]interface{}{}
		err = json.Unmarshal([]byte(orderInDB.ChangeLog), &currentChangeLog)
		if err != nil {
			return fmt.Errorf("%w: error al deserializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Agregar order al arreglo
		newRecord := map[string]interface{}{}
		err = json.Unmarshal(orderInDB.ToJson(), &newRecord)
		if err != nil {
			return fmt.Errorf("%w: error covertir orden en json: %v", errorhandling.ErrInternal, err)
		}
		// Remove irrelevant data
		delete(newRecord, "change_log")
		delete(newRecord, "id")
		delete(newRecord, "author_id")
		delete(newRecord, "created_at")

		// Add updater name
		author, err := s.userRepo.ReadByID(orderInDB.UpdatedBy)
		if err != nil {
			return fmt.Errorf("%w: error al obtener autor: %v", errorhandling.ErrInternal, err)
		}
		newRecord["updated_by_name"] = fmt.Sprintf("%s (%s)", author.UserFullName, author.Email)
		currentChangeLog = append(currentChangeLog, newRecord)

		// Convertir de nuevo a string
		newChangeLog, err := json.MarshalIndent(currentChangeLog, "", "  ")
		if err != nil {
			return fmt.Errorf("%w: error al serializar change log: %v", errorhandling.ErrInternal, err)
		}

		// Actualizar change log
		updates.ChangeLog = datatypes.JSON(newChangeLog)
	}

	// Actualizar
	return s.repo.Update(updates)
}

func (s *orderService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}
