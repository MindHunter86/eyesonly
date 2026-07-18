package sessions

import (
	"time"

	"github.com/MindHunter86/eyesonly/internal/webapp/shared/httpx"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(router fiber.Router) {
	router.Post("/session", h.Create)
}

func (h *Handler) Create(c *fiber.Ctx) (e error) {
	var token *Token
	if token, e = h.svc.Issue(time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, fiber.Map{
		"token":      token.Token,
		"session":    token.Session,
		"expires_at": token.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
