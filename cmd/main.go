package main

import (
	"fmt"
	"time"

	"github.com/Abhishek2010dev/kokoro"
)

func main() {
	server := kokoro.New()

	server.Use(func(ctx *kokoro.Context, next kokoro.HandlerFunc) error {
		start := time.Now()
		err := next(ctx)
		fmt.Printf("[%s] %s  (%s)\n", ctx.Method(), ctx.Path(), time.Since(start))
		return err
	})

	// Main route: /
	server.GET("/", func(ctx *kokoro.Context) error {
		data := map[string]any{
			"message": "Hello, World",
			"status":  true,
			"count":   123,
		}
		return ctx.JSON(data)
	})

	// Optional: Route group for /api
	api := server.Group("/api")
	api.GET("/status", func(ctx *kokoro.Context) error {
		return ctx.JSON(map[string]string{"status": "ok"})
	})

	// Start the server
	server.ListenAndServe(":3000")
}
