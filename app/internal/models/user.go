package models

import (
	"time"

	"gorm.io/datatypes"
)

// User representa la estructura de un usuario en el sistema
type User struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Username contiene el nombre de usuario asignado por el usuario
	UserFullName string `gorm:"type:varchar(100);not null;index:idx_user_full_name" json:"user_full_name"`

	// Email contiene el correo electronico asignado por el usuario
	Email string `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`

	// Phone contiene el numero de telefono asignado por el usuario
	Phone string `gorm:"type:varchar(20)" json:"phone"`

	// Password contiene la contraseña recibida del usuario
	Password string `gorm:"-" json:"password"`

	// PasswordHash contiene la contraseña encriptada del usuario
	PasswordHash string `gorm:"type:varchar(255);not null" json:"password_hash"`

	// RequiresPasswordUpdate indica si el password debe ser actualizado por el usuario durante un inicio de sesion
	RequiresPasswordUpdate bool `gorm:"not null;default:false" json:"requires_password_update"`

	// Rol definido por el usuario
	Role UserRole `gorm:"type:varchar(20);not null;default:'User';index:idx_user_role" json:"role"`

	// Estado definido por el usuario
	Status UserStatus `gorm:"type:varchar(20);not null;default:'Active';index:idx_user_status" json:"status"`

	// Fechas y auditoría asignadas automáticamente por el sistema
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// ChangeHistory registra el historial de cambios
	ChangeHistory datatypes.JSON `gorm:"type:jsonb" json:"change_history"`

	// StoredInChangeLogAt registra la fecha en la que se almaceno el cambio
	StoredInChangeLogAt time.Time `json:"stored_in_change_log_at"`
}

// TableName especifica explícitamente el nombre de la tabla en singular para GORM.
func (User) TableName() string {
	return "user"
}
