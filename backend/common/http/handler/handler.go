package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/qosdil/like-x/backend/common/service"
)

func ErrResp(c fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}

func HandleCommonErrs(c fiber.Ctx, err error) error {
	switch err {
	case service.ErrAlreadyExists:
		return c.SendStatus(http.StatusConflict)
	case service.ErrBadRequest:
		return ErrResp(c, fiber.StatusBadRequest, "bad_request", "Bad request")
	case service.ErrForbidden:
		return c.SendStatus(fiber.StatusForbidden)
	case service.ErrNotFound:
		return ErrResp(c, fiber.StatusNotFound, "record_not_found", "Record not found")
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func ObjResp(c fiber.Ctx, data any, err error) error {
	if err != nil {
		return HandleCommonErrs(c, err)
	}

	return c.JSON(data)
}
