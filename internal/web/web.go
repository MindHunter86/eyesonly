package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// Dist contains the built frontend. The directory is intentionally inside this
// Go package because go:embed cannot embed files from parent directories.
//
//go:embed all:dist
var Dist embed.FS

func RegisterStatic(app *fiber.App) error {
	dist, err := fs.Sub(Dist, "dist")
	if err != nil {
		return err
	}

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if !strings.HasPrefix(c.Path(), "/v1/") {
			c.Set("Cache-Control", "no-store")
		}
		return err
	})

	app.Use("/", filesystem.New(filesystem.Config{Root: http.FS(dist), Browse: false, Index: "index.html", NotFoundFile: "index.html"}))
	return nil
}
