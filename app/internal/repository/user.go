package repository

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/models/filters"
	"fmt"

	"gorm.io/gorm"
)

// Interfaces para permitir testing
type UserRepository interface {
	Create(user models.User) error
	ReadByFilter(filters filters.GetUsers) ([]models.User, int64, error)
	ReadByQuery(query models.User) ([]models.User, error)
	ReadByID(id uint) (models.User, error)
	Update(user models.User) error
	Delete(id uint) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Implementación de User
type userRepo struct{ db *gorm.DB }

func (r *userRepo) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepo) ReadByFilter(filters filters.GetUsers) ([]models.User, int64, error) {
	// Buils query
	var users []models.User

	query := r.db.Debug().Model(&models.User{})

	if filters.UserFullName != "" {
		query = query.Where("immutable_unaccent(\"user\".user_full_name) ILIKE immutable_unaccent(?)", "%"+filters.UserFullName+"%")
	}

	if filters.Role != "" {
		query = query.Where("\"user\".role = ?", filters.Role)
	}

	if filters.Email != "" {
		query = query.Where("\"user\".email ILIKE ?", "%"+filters.Email+"%")
	}

	if filters.Status != "" && filters.Status != "all" {
		query = query.Where("\"user\".status = ?", filters.Status)
	}

	// Get counter
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return users, 0, fmt.Errorf("%w: error al contar registros", errorhandling.ErrInternal)
	}

	// 5. Aplicamos Paginación (LIMIT y OFFSET) y traemos los datos reales
	offset := (filters.Page - 1) * filters.Limit
	err := query.Limit(filters.Limit).Offset(offset).Order("\"user\".created_at DESC").Find(&users).Error
	return users, totalRecords, err
}

func (r *userRepo) ReadByQuery(query models.User) ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users, query).Error
	return users, err
}

func (r *userRepo) ReadByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return models.User{}, errorhandling.ErrNotFoundError
	}

	return user, err
}

func (r *userRepo) Update(user models.User) error {
	return r.db.Debug().Updates(&user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
