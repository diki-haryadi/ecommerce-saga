package utils

type Pagination struct {
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
	Total    int64 `json:"total,omitempty"`
}

func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	return p.PageSize
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}

type PaginatedResult struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginatedResult(items interface{}, pagination *Pagination) *PaginatedResult {
	totalPages := int(pagination.Total) / pagination.PageSize
	if int(pagination.Total)%pagination.PageSize > 0 {
		totalPages++
	}

	return &PaginatedResult{
		Items:      items,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: pagination.Total,
		TotalPages: totalPages,
	}
}
