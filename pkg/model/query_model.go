package model

type Query struct {
	// Date range for filtering
	FromDate string `query:"from_date"`
	ToDate   string `query:"to_date"`

	// Pagination parameters
	Page     int `query:"page"`
	PageSize int `query:"page_size"`

	// Sorting parameter
	SortBy string `query:"sort_by"`

	// Sub-struct expansion
	Expand []string `query:"expand"`

	// Custom query parameters
	Custom map[string]interface{}
}
