package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"subscriber/router"
	"subscriber/utils"
)

func main() {
	r := router.SetupRouter()

	log.Fatalf("gin stopped!%s", r.Run(":8080"))
}
