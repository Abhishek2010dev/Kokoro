package main

import "github.com/Abhishek2010dev/kokoro"

func main() {
	server := kokoro.New()
	server.GET("/", func(ctx *kokoro.Context) error {
		return ctx.Json(map[string]any{"messsage": "hello"})
	})
	server.ListenAndServe(":3000")
}
