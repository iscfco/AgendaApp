package models

type UserRole string

const (
	UserRoleUser       UserRole = "user"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSuperAdmin UserRole = "superadmin"
)

type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusDisabled UserStatus = "disabled"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusAll       OrderStatus = "all"
)

var OrderStatusText map[OrderStatus]string = map[OrderStatus]string{
	OrderStatusPending:   "Pendiente",
	OrderStatusDelivered: "Entregado",
	OrderStatusAll:       "Todos",
}
