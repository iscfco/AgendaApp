package models

import "time"

// User representa la estructura de un usuario en el sistema [cite: 5, 6]
type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Username contiene el nombre de usuario asignado por el usuario
	Username string `gorm:"unique;not null" json:"username"`

	// Email contiene el correo electronico asignado por el usuario
	Email string `gorm:"unique;not null" json:"email"`

	// Phone contiene el numero de telefono asignado por el usuario
	Phone string `json:"phone"`

	// Password suele ser procesada por el sistema (hashing) antes de guardarse
	Password string `gorm:"not null" json:"-"`

	// UpdatePassword indica si el password debe ser actualizado por el usuario durante un inicio de sesion
	UpdatePassword bool `json:"update_password"`

	// Rol definido por el usuario
	Role UserRole `json:"role"`

	// Estado definido por el usuario
	Status UserStatus `json:"status"`

	// Fechas y auditoría asignadas automáticamente por el sistema
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// ChangeLog registra el historial de cambios
	ChangeLog string `gorm:"type:json" json:"change_log"`

	// StoredInChangeLogAt registra la fecha en la que se almaceno el cambio
	StoredInChangeLogAt time.Time `json:"stored_in_change_log_at"`
}
