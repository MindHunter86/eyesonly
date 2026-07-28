package secrets

import (
	"context"
	"strings"
	"time"

	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/MindHunter86/eyesonly/internal/webapp/sessions"
	"github.com/MindHunter86/eyesonly/internal/webapp/shared/apperr"
	"github.com/MindHunter86/eyesonly/internal/webapp/shared/httpx"
	"github.com/urfave/cli/v2"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	sessions *sessions.Service
	secrets  *Service

	cli *cli.Context
}

func NewHandler(c context.Context, sessions *sessions.Service, secrets *Service) *Handler {
	cli := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)

	return &Handler{
		sessions: sessions,
		secrets:  secrets,

		cli: cli,
	}
}

func (h *Handler) Register(router fiber.Router) {
	router.Get("/stats", h.Stats)

	router.Get("/session/secrets", h.ListSession)

	router.Post("/secrets", h.Create)
	router.Post("/destroy", h.DestroyByToken)

	router.Get("/secrets/:id", h.Get)
	router.Post("/secrets/:id/reveal", h.Reveal)
	router.Post("/secrets/:id/destroy", h.DestroyBySession)

	router.Get("/admin/secrets", h.requireAdmin, h.ListAdmin)
}

func (h *Handler) Create(c *fiber.Ctx) (e error) {
	var sess *sessions.Session
	if sess, e = h.sessionFromRequest(c); e != nil {
		return
	}

	var input CreateInput
	if e = c.BodyParser(&input); e != nil {
		return apperr.Validation("Invalid request body")
	}

	var created *Created
	if created, e = h.secrets.Create(c.Context(), sess, input, time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, created)
}

func (h *Handler) ListSession(c *fiber.Ctx) (e error) {
	var sess *sessions.Session
	if sess, e = h.sessionFromRequest(c); e != nil {
		return
	}

	var items []Public
	if items, e = h.secrets.ListSession(c.Context(), sess, time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, fiber.Map{"items": items})
}

func (h *Handler) Get(c *fiber.Ctx) (e error) {
	var item *Public
	if item, e = h.secrets.GetPublic(c.Context(), c.Params("id"), time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, item)
}

func (h *Handler) Reveal(c *fiber.Ctx) (e error) {
	var item *Reveal
	if item, e = h.secrets.Reveal(c.Context(), c.Params("id"), time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, item)
}

func (h *Handler) DestroyByToken(c *fiber.Ctx) (e error) {
	var input struct {
		DestroyToken string `json:"destroy_token"`
	}

	if e = c.BodyParser(&input); e != nil {
		return apperr.InvalidDestroyToken()
	}

	var result *DestroyResult
	if result, e = h.secrets.DestroyByToken(c.Context(), input.DestroyToken, time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, result)
}

func (h *Handler) DestroyBySession(c *fiber.Ctx) (e error) {
	var sess *sessions.Session
	if sess, e = h.sessionFromRequest(c); e != nil {
		return
	}

	var result *DestroyResult
	if result, e = h.secrets.DestroyBySession(c.Context(), sess, c.Params("id"), time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, result)
}

func (h *Handler) Stats(c *fiber.Ctx) (e error) {
	var stats *Stats
	if stats, e = h.secrets.Stats(c.Context(), time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, stats)
}

func (h *Handler) ListAdmin(c *fiber.Ctx) (e error) {
	filter := &AdminFilter{
		Q:        c.Query("q"),
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 10),
	}

	var items *AdminList
	if items, e = h.secrets.ListAdmin(c.Context(), filter, time.Now().UTC()); e != nil {
		return
	}

	return httpx.OK(c, items)
}

func (h *Handler) sessionFromRequest(c *fiber.Ctx) (sess *sessions.Session, e error) {
	var token string

	if token = strings.TrimSpace(c.Get("Authorization")); token == "" {
		token = strings.TrimSpace(c.Get("X-EyesOnly-Session"))
	}

	return h.sessions.Parse(token)
}

func (h *Handler) requireAdmin(c *fiber.Ctx) error {
	if h.cli.String("admin-access-token") == "" {
		return c.Next()
	}

	var token string
	if token = strings.TrimSpace(c.Get("X-EyesOnly-Admin-Token")); token == "" {
		token = strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	}

	if token != h.cli.String("admin-access-token") {
		return apperr.Unauthorized(apperr.CodeAdminRequired, "Admin token is required.")
	}

	return c.Next()
}
