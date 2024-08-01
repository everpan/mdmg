package main

import (
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
)

func main() {
	app := fiber.New()
	app.Get("/spike", func(c *fiber.Ctx) error {
		v8Ctx := v8.NewContext()
		defer v8Ctx.Close()

		// c.Request().Header.co
		return nil
	})
	app.Listen(":9090")
}
