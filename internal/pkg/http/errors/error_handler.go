package errors

import (
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
	"github.com/gofiber/fiber/v2"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ValidationError     ErrorType = "VALIDATION_ERROR"
	AuthenticationError ErrorType = "AUTHENTICATION_ERROR"
	NotFoundError       ErrorType = "NOT_FOUND_ERROR"
	ConflictError       ErrorType = "CONFLICT_ERROR"
	InternalError       ErrorType = "INTERNAL_ERROR"
)

// AppError represents application specific error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// ErrorHandler interface for Chain of Responsibility pattern
type ErrorHandler interface {
	Handle(c *fiber.Ctx, err error) error
	SetNext(handler ErrorHandler) ErrorHandler
}

// BaseHandler provides base implementation
type BaseHandler struct {
	next ErrorHandler
}

func (h *BaseHandler) SetNext(handler ErrorHandler) ErrorHandler {
	h.next = handler
	return handler
}

func (h *BaseHandler) Handle(c *fiber.Ctx, err error) error {
	if h.next != nil {
		return h.next.Handle(c, err)
	}
	return response.InternalServerError(c, "An unexpected error occurred")
}

// ValidationErrorHandler handles validation errors
type ValidationErrorHandler struct {
	BaseHandler
}

func (h *ValidationErrorHandler) Handle(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok && appErr.Type == ValidationError {
		return response.BadRequest(c, appErr.Message)
	}
	return h.BaseHandler.Handle(c, err)
}

// AuthenticationErrorHandler handles authentication errors
type AuthenticationErrorHandler struct {
	BaseHandler
}

func (h *AuthenticationErrorHandler) Handle(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok && appErr.Type == AuthenticationError {
		return response.Unauthorized(c, appErr.Message)
	}
	return h.BaseHandler.Handle(c, err)
}

// NotFoundErrorHandler handles not found errors
type NotFoundErrorHandler struct {
	BaseHandler
}

func (h *NotFoundErrorHandler) Handle(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok && appErr.Type == NotFoundError {
		return response.NotFound(c, appErr.Message)
	}
	return h.BaseHandler.Handle(c, err)
}

// ConflictErrorHandler handles conflict errors
type ConflictErrorHandler struct {
	BaseHandler
}

func (h *ConflictErrorHandler) Handle(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok && appErr.Type == ConflictError {
		return response.Error(c, fiber.StatusConflict, appErr.Message)
	}
	return h.BaseHandler.Handle(c, err)
}

// NewErrorHandler creates a chain of error handlers
func NewErrorHandler() ErrorHandler {
	validationHandler := &ValidationErrorHandler{}
	authHandler := &AuthenticationErrorHandler{}
	notFoundHandler := &NotFoundErrorHandler{}
	conflictHandler := &ConflictErrorHandler{}

	validationHandler.SetNext(authHandler).
		SetNext(notFoundHandler).
		SetNext(conflictHandler)

	return validationHandler
}

// Error creation helper functions
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: message,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Type:    AuthenticationError,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ConflictError,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Type:    InternalError,
		Message: "Internal server error",
		Err:     err,
	}
}
