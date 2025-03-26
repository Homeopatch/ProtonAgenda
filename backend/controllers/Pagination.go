package controllers

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page" doc:"The current page number." example:"1"`
	PageSize   int `json:"pageSize" doc:"The number of items per page." example:"20"`
	TotalItems int `json:"totalItems" doc:"The total number of items available." example:"123"`
	TotalPages int `json:"totalPages" doc:"The total number of pages available." example:"7"`
}
