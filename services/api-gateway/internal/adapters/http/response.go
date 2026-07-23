package http

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func Paginated(c *gin.Context, status int, message string, data interface{}, pagination Pagination) {
	c.JSON(status, PaginatedResponse{
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Error(c *gin.Context, status int, error error) {
	c.JSON(status, ErrorResponse{
		Error: error.Error(),
	})
}
