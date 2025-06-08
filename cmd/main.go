package main

import "github.com/Abhishek2010dev/kokoro"

func main() {
	server := kokoro.New()
	server.GET("/{name}", func(ctx *kokoro.Context) error {
		return ctx.Text("Hello, " + ctx.Param("name"))
	})
	server.ListenAndServe(":3000")
}
