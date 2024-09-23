package main

import (
	"gin_server/router"
	"log"
)

func main() {
	r := router.SetupRouter()

	log.Fatalf("gin stopped!%s", r.Run(":8080"))
}
