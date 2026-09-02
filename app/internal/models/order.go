package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Order representa la estructura de un pedido en el sistema
type Order struct {
	// Author es el usuario que creo el pedido
	Author User `gorm:"foreignKey:AuthorID" json:"-"`

	// ID es el identificador del pedido
	ID uint `gorm:"primaryKey" json:"id"`

	// Autor es el usuario que creo el pedido
	AuthorID uint `gorm:"not null;foreignKey" json:"author_id"`

	// ClienteName es el nombre completo del cliente
	ClientName string `gorm:"not null" json:"client_name"`

	// ClientPhone es el número de teléfono del cliente
	ClientPhone string `json:"client_phone"`

	// ClientAddress es la dirección del cliente
	ClientAddress string `json:"client_address"`

	// TotalPrice es el precio total del pedido
	TotalPrice float64 `json:"total_price"`

	// DownPayment es el anticipo del pedido
	DownPayment float64 `json:"down_payment"`

	// CreatedAt es la fecha y hora de creación del pedido
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt es la fecha y hora de la última actualización del pedido
	UpdatedAt time.Time `json:"updated_at"`

	// UpdatedBy es el ID del usuario que actualizo el pedido
	UpdatedBy uint `json:"updated_by"`

	// DeliveryDate es la fecha de entrega del pedido
	DeliveryDate time.Time `json:"delivery_date"`

	// Description es la descripción del pedido
	Description string `json:"description"`

	// Estado del pedido
	Status OrderStatus `json:"status"`

	// ChangeLog es el historial de cambios del pedido
	ChangeLog datatypes.JSON `json:"change_log"`
}

func (Order) TableName() string {
	return "order"
}

func (o *Order) ToJson() []byte {
	bytes, err := json.Marshal(o)
	if err != nil {
		return []byte("")
	}
	return bytes
}
