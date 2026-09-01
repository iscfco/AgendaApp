package bootstrap

import (
	"fmt"

	"gorm.io/gorm"

	"agenda-app/app/internal/errorhandling"
	"agenda-app/app/internal/models"
	"agenda-app/app/internal/utils"
)

func Run(db *gorm.DB, email, pwd, name string) error {
	if email == "" {
		return fmt.Errorf("bootstrap admin email is required")
	}

	if pwd == "" {
		return fmt.Errorf("bootstrap admin password is required")
	}

	var user models.User

	err := db.
		Where("email = ?", email).
		First(&user).Error

	if err == nil {
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("check bootstrap admin: %w", err)
	}

	passwordHash, err := utils.HashPassword(pwd)
	if err != nil {
		return fmt.Errorf("%w: error al encryptar password: %v", errorhandling.ErrInternal)
	}

	user = models.User{
		UserFullName:           name,
		Email:                  email,
		PasswordHash:           string(passwordHash),
		RequiresPasswordUpdate: false, // TODO: change for true
		Role:                   models.UserRoleSuperAdmin,
		Status:                 models.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("create bootstrap admin: %w", err)
	}

	return nil
}
