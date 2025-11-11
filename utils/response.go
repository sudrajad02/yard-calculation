package utils

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseFormat struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func ApiResponse(c *fiber.Ctx, code int, message string, data any, error any) {
	jsonResponse := ResponseFormat{
		Code:    code,
		Message: message,
	}

	if data != nil {
		jsonResponse.Data = data
	}

	if error != nil {
		jsonResponse.Error = error
	}

	c.Status(code).JSON(jsonResponse)
}
