package controllers

import "github.com/gofiber/fiber/v2"

func NewApiV1(root *fiber.Router) {
	apiv1 := *root

	// customer service API
	adm := apiv1.Group("/admin")
	adm.Get("/whoami")

	// account routes
	acc := apiv1.Group("/account")

	acc.Post("/login")
	acc.Post("/login/totp")
	acc.Post("/login/get-token")

	acc.Post("/register")
	acc.Post("/register/verify")

	acc.Post("/logout")

	// proxy subject routes
	prx := apiv1.Post("/proxy")

	prx.Post("/new")

	prx.Post("/:id/info")
	prx.Post("/:id/enable")
	prx.Post("/:id/disable")
	prx.Post("/:id/remove")

	// internal routes
	apiv1.Get("/status")
	apiv1.Get("/version")
}

/*
- JWT
- store

- secrets
- sessions ?
- stats
- admin


*/
