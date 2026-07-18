package auth

import (
	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func (m *Handler) Register(c *fiber.Ctx) error {
	//

	return utils.AcquireFiberError(1, "")
}

func (m *Handler) Login(c *fiber.Ctx) error { return utils.AcquireFiberError(1, "") }

func (m *Handler) Refresh(c *fiber.Ctx) error { return utils.AcquireFiberError(1, "") }

func (m *Handler) Logout(c *fiber.Ctx) error { return utils.AcquireFiberError(1, "") }

func (m *Handler) LogoutAll(c *fiber.Ctx) error { return utils.AcquireFiberError(1, "") }
