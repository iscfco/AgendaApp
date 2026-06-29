package repository

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order models.Order) error
	ReadByQuery(order models.Order) ([]models.Order, error)
	ReadByID(id uint) (models.Order, error)
	Update(order models.Order) error
	Delete(id uint) error
}

func NewOrderFactory(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

type orderRepo struct{ db *gorm.DB }

func (r *orderRepo) Create(o models.Order) error {
	return r.db.Create(o).Error
}

func (r *orderRepo) ReadByQuery(order models.Order) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Find(&orders, order).Error
	return orders, err
}

func (r *orderRepo) ReadByID(id uint) (models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	if err == gorm.ErrRecordNotFound {
		return models.Order{}, errorhandling.ErrNotFoundError
	}

	return order, err
}

func (r *orderRepo) Update(order models.Order) error {
	return r.db.Save(&order).Error
}

func (r *orderRepo) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}
