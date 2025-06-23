package response

import (
	"github.com/gofiber/fiber/v2"
)

// Standard response structure
type Response struct {
	Status  int         `json:"-"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseBuilder implements Template Method pattern
type ResponseBuilder interface {
	SetStatus(status int) ResponseBuilder
	SetSuccess(success bool) ResponseBuilder
	SetMessage(message string) ResponseBuilder
	SetError(err string) ResponseBuilder
	SetData(data interface{}) ResponseBuilder
	Build() *Response
}

type responseBuilder struct {
	response *Response
}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder() ResponseBuilder {
	return &responseBuilder{
		response: &Response{},
	}
}

func (rb *responseBuilder) SetStatus(status int) ResponseBuilder {
	rb.response.Status = status
	return rb
}

func (rb *responseBuilder) SetSuccess(success bool) ResponseBuilder {
	rb.response.Success = success
	return rb
}

func (rb *responseBuilder) SetMessage(message string) ResponseBuilder {
	rb.response.Message = message
	return rb
}

func (rb *responseBuilder) SetError(err string) ResponseBuilder {
	rb.response.Error = err
	return rb
}

func (rb *responseBuilder) SetData(data interface{}) ResponseBuilder {
	rb.response.Data = data
	return rb
}

func (rb *responseBuilder) Build() *Response {
	return rb.response
}

// Factory methods for common responses
func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := NewResponseBuilder().
		SetStatus(status).
		SetSuccess(true).
		SetMessage(message).
		SetData(data).
		Build()

	return c.Status(response.Status).JSON(response)
}

func Error(c *fiber.Ctx, status int, message string) error {
	response := NewResponseBuilder().
		SetStatus(status).
		SetSuccess(false).
		SetError(message).
		Build()

	return c.Status(response.Status).JSON(response)
}

// Common HTTP status code responses
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

func InternalServerError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return Error(c, fiber.StatusInternalServerError, message)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusCreated, message, data)
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusOK, message, data)
}
