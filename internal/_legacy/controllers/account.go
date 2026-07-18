package controllers

import (
	"errors"

	"github.com/MindHunter86/eyesonly/internal/models"
	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type account struct {
	Email    []byte
	Password []byte
}

func AccountLogin() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (e error) {
		// add some headers validators

		// parse email and password
		reqacc := new(account)
		if e = c.BodyParser(&reqacc); e != nil {
			return utils.AcquireFiberError(fiber.StatusInternalServerError, e.Error())
		}

		ctx := c.UserContext()

		var acc *models.Account
		if acc, e = models.LoginAccount(ctx, reqacc.Email, reqacc.Password); e != nil {
			if errors.Is(e, models.ErrNoAccFound) {
				return utils.AcquireFiberError(fiber.StatusForbidden, e.Error())
			} else if errors.Is(e, models.ErrAccIsBlocked) {
				return utils.AcquireFiberError(fiber.StatusForbidden, e.Error())
			} else {
				return utils.AcquireFiberError(fiber.StatusInternalServerError, e.Error())
			}
		}

		// todo ; migrate to easyjson
		return c.JSON(acc)
	}
}

func AccountRegister() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (e error) {
		// add some headers validators

		// parse email and password
		reqacc := new(account)
		if e = c.BodyParser(&reqacc); e != nil {
			return utils.AcquireFiberError(fiber.StatusInternalServerError, e.Error())
		}

		ctx := c.UserContext()

		var acc *models.Account
		if acc, e = models.CreateAccount(ctx, reqacc.Email, reqacc.Password); e != nil {
			return utils.AcquireFiberError(fiber.StatusInternalServerError, e.Error())
		}

		return c.JSON(acc)
	}
}
