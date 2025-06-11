package main

import (
	"github.com/Abhishek2010dev/kokoro"
)

type Message struct {
	XMLName string `xml:"message"`
	Text    string `xml:",chardata"`
}

func main() {
	server := kokoro.New()

	server.GET("/", func(ctx *kokoro.Context) error {
		msg := Message{Text: "Hello, XML"}
		return ctx.XML(msg)
	})

	server.ListenAndServe(":3000")
}
