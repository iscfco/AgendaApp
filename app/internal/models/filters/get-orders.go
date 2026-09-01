package filters

type GetOrders struct {
	Show            bool   `form:"show"`
	UserCreatorName string `form:"user_creator_name"`
	Keyword         string `form:"keyword"`
	ClientName      string `form:"client_name"`
	From            string `form:"from_date"`
	To              string `form:"to_date"`
	Status          string `form:"status"`

	// Pagination fields
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
