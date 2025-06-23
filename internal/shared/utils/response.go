package utils

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page      int   `json:"page,omitempty"`
	PageSize  int   `json:"page_size,omitempty"`
	TotalRows int64 `json:"total_rows,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(code, message string) Response {
	return Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}

func PaginatedResponse(data interface{}, page, pageSize int, totalRows int64) Response {
	return Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:      page,
			PageSize:  pageSize,
			TotalRows: totalRows,
		},
	}
}
