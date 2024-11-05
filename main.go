package main

import (
	"CloudBook/internal/web"
)

func main() {
	server := web.RegisterRoutes()
	server.Run(":8080")
}
