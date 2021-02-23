package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	funcs "github.com/ONSDigital/blaise-nifi-encrypt"
)

// emulates the cloud functions
func main() {
	funcframework.RegisterEventFunction("/encrypt", funcs.EncryptFunction)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
