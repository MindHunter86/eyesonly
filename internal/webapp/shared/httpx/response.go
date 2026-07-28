package httpx

import (
	"errors"

	"github.com/MindHunter86/eyesonly/internal/webapp/shared/apperr"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	OK    bool       `json:"ok"`
	Data  any        `json:"data,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		return c.Status(appErr.Status).JSON(Response{OK: false, Error: &ErrorBody{Code: string(appErr.Code), Message: appErr.Message, Details: appErr.Details}})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(Response{OK: false, Error: &ErrorBody{Code: "HTTP_ERROR", Message: fiberErr.Message}})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(Response{OK: false, Error: &ErrorBody{Code: string(apperr.CodeInternal), Message: "Internal server error."}})
}

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(Response{OK: true, Data: data})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{OK: true, Data: data})
}
