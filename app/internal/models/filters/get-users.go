package filters

type GetUsers struct {
	Show         bool   `form:"show" json:"show"`
	UserFullName string `form:"user_full_name" json:"user_full_name"`
	Role         string `form:"role" json:"role"`
	Email        string `form:"email" json:"email"`
	Status       string `form:"status" json:"status"`

	// Pagination fields
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}
