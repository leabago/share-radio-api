package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Tasks(ctx *fiber.Ctx) error {

	fmt.Println("Tasks middleware")

	return ctx.Next()
}
