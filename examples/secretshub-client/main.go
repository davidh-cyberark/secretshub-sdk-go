package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
	"github.com/davidh-cyberark/secretshub-sdk-go/secretshub"
)

var (
	version = "dev"
)

func main() {
	idtenanturl := flag.String("idtenanturl", "", "Identity URL")
	iduser := flag.String("iduser", "", "Identity user id")
	idpass := flag.String("idpass", "", "Identity user password")
	shurl := flag.String("shurl", "", "Secrets Hub URL, Ex: https://EXAMPLE.secretshub.cyberark.cloud/")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	// logger
	logger := log.New(os.Stderr, "[secretshub-client] ", log.LstdFlags)

	if !*debug {
		logger.SetOutput(io.Discard)
	}

	ctx := context.Background()

	// Create the Authentication provider
	userAuth := &identity.UserCredentialsAuthenticationProvider{
		User: *iduser,
		Pass: *idpass,
	}

	// Create the Identity client with the authentication provider
	idClient, idClientErr := identity.NewClientWithResponses(*idtenanturl,
		identity.WithRequestEditorFn(userAuth.Intercept))
	if idClientErr != nil {
		logger.Fatalf("failed to create id client: %v", idClientErr)
	}

	// Create the Identity service with the client and authentication provider
	service := &identity.Service{
		TenantURL:     *idtenanturl,
		Client:        idClient,
		Logger:        logger,
		AuthnProvider: userAuth,
	}

	ctx = context.WithValue(ctx, identity.IdentityService, service)

	client, clientErr := secretshub.NewClientWithResponses(*shurl,
		secretshub.WithRequestEditorFn(userAuth.Intercept))
	if clientErr != nil {
		logger.Fatalf("failed to create secretshub client: %v", clientErr)
	}

	if client == nil {
		logger.Fatalf("failed to create secretshub client")
	}

	logger.Println("Successfully created Secrets Hub client.")

	// Use the client to interact with the Secrets Hub API
}
