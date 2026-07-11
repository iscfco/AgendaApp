package repository

import (
	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"

	"gorm.io/gorm"
)

// Interfaces para permitir testing
type UserRepository interface {
	Create(user models.User) error
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
	return r.db.Create(user).Error
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
	return r.db.Save(&user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
