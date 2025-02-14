package main

import (
	"github.com/bryanvaz/go-lucide/test/pages"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return pages.Index().Render(c.Context(), c.Response().BodyWriter())
	})

	app.Listen(":3000")
}
