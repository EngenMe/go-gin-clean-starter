package dto

type (
	// PaginationRequest represents a request structure for managing pagination parameters such as search, page, and per page.
	PaginationRequest struct {
		Search  string `form:"search"`
		Page    int    `form:"page"`
		PerPage int    `form:"per_page"`
	}

	// PaginationResponse represents metadata for paginated responses including page, items per page, total pages, and item count.
	PaginationResponse struct {
		Page    int   `json:"page"`
		PerPage int   `json:"per_page"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)

// GetOffset calculates the offset for pagination based on the current page and records per page.
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the number of records to be retrieved per page as defined in the PerPage field of the request.
func (p *PaginationRequest) GetLimit() int {
	return p.PerPage
}

// GetPage retrieves the current page number from the PaginationRequest structure.
func (p *PaginationRequest) GetPage() int {
	return p.Page
}

// Default sets default values for Page and PerPage if they are not provided (zero values).
func (p *PaginationRequest) Default() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PerPage == 0 {
		p.PerPage = 10
	}
}
