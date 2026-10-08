package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daytona/clients/sdk-go/pkg/daytona"
	modal "github.com/modal-labs/modal-client/go"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "daytona":
		mustDaytona()
	case "modal":
		mustModal()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`sandbox-providers — Daytona & Modal Go SDKs

Usage:
  sandbox-providers daytona   Create Daytona client (DAYTONA_API_KEY or DAYTONA_JWT_TOKEN)
  sandbox-providers modal     Create Modal client (~/.modal.toml or MODAL_TOKEN_ID/SECRET)

Examples:
  go run . daytona
  go run . modal`)
}

func mustDaytona() {
	client, err := daytona.NewClient()
	if err != nil {
		log.Fatalf("daytona: %v", err)
	}
	fmt.Println("Daytona client ready")
	_ = client
}

func mustModal() {
	client, err := modal.NewClient()
	if err != nil {
		log.Fatalf("modal: %v", err)
	}
	defer client.Close()
	fmt.Printf("Modal client ready (SDK %s)\n", client.Version())
}
