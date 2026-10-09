package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Trả về 200 OK
func OK[T any](c *gin.Context, data T, message string) {
	c.JSON(http.StatusOK, Success(data, message))
}

// Trả về 201 Created
func Created[T any](c *gin.Context, data T, message string) {
	c.JSON(http.StatusCreated, Success(data, message))
}

// Trả về lỗi chung với statusCode tuỳ biến
func Error(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, ErrorResponse(message, err))
}

// Trả về 400 Bad Request
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message, err)
}

// Trả về 401 Unauthorized
func Unauthorized(c *gin.Context, message string, err error) {
	Error(c, http.StatusUnauthorized, message, err)
}

// Trả về 404 Not Found
func NotFound(c *gin.Context, message string, err error) {
	Error(c, http.StatusNotFound, message, err)
}

// Trả về 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}
