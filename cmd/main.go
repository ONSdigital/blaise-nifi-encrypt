package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	funcs "github.com/ONSDigital/blaise-mi-extract"
	"log"
	"os"
)

// emulates the cloud functions
func main() {

	funcframework.RegisterEventFunction("/extract", funcs.ExtractFunction)
	funcframework.RegisterEventFunction("/zip", funcs.ZipFunction)
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
