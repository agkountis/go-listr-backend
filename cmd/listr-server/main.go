package main

import (
	"fmt"
	"log"
	"os"

	"github.com/agkountis/go-listr-backend/internal/app/listr-server/server"
)

var selfSignedCertsPath string = os.Getenv("DOMAIN_SELF_SIGNED_CERTS_PATH")
var certFilePath string = fmt.Sprintf("%v/localhost.pem", selfSignedCertsPath)
var keyFilePath string = fmt.Sprintf("%v/localhost-key.pem", selfSignedCertsPath)

func main() {
	server, err := server.New()

	if err != nil {
		log.Fatalf("Failed to create server. Error: %v", err)
	}

	server.StartTLS("0.0.0.0:8080", certFilePath, keyFilePath)
}
