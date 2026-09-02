package repository

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/models/filters"
	"fmt"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order models.Order) error
	ReadByQuery(filter filters.GetOrders) ([]models.Order, int64, error)
	ReadByID(id uint) (models.Order, error)
	Update(order models.Order) error
	Delete(id uint) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

type orderRepo struct{ db *gorm.DB }

func (r *orderRepo) Create(o models.Order) error {
	return r.db.Create(&o).Error
}

func (r *orderRepo) ReadByQuery(filters filters.GetOrders) ([]models.Order, int64, error) {
	// Buils query
	var orders []models.Order

	query := r.db.Model(&models.Order{}).Joins("Author")

	if filters.UserCreatorName != "" {
		query = query.Where("immutable_unaccent(\"Author\".user_full_name) ILIKE immutable_unaccent(?)", "%"+filters.UserCreatorName+"%")
	}

	if filters.ClientName != "" {
		query = query.Where("immutable_unaccent(\"order\".client_name) ILIKE immutable_unaccent(?)", "%"+filters.ClientName+"%")
	}

	if filters.Keyword != "" {
		query = query.Where("immutable_unaccent(\"order\".description) ILIKE immutable_unaccent(?)", "%"+filters.Keyword+"%")
	}

	if filters.From != "" {
		query = query.Where("\"order\".created_at >= ?", filters.From)
	}

	if filters.To != "" {
		query = query.Where("\"order\".created_at <= ?", filters.To)
	}

	if filters.Status != "" && filters.Status != "all" {
		query = query.Where("\"order\".status = ?", filters.Status)
	}

	// Get counter
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return orders, 0, fmt.Errorf("%w: error al contar registros", errorhandling.ErrInternal)
	}

	// 5. Aplicamos Paginación (LIMIT y OFFSET) y traemos los datos reales
	offset := (filters.Page - 1) * filters.Limit
	err := query.Limit(filters.Limit).Offset(offset).Find(&orders).Error
	return orders, totalRecords, err
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
