package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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

	slogOpts := &slog.HandlerOptions{}
	if *debug {
		slogOpts.Level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, slogOpts).WithAttrs([]slog.Attr{
		slog.String("service", "secretshub"),
		slog.String("command", "secretshub-client"),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
		slog.Error(fmt.Sprintf("failed to create id client: %v", idClientErr))
		os.Exit(1)
	}

	// Create the Identity service with the client and authentication provider
	service := &identity.Service{
		TenantURL:     *idtenanturl,
		Client:        idClient,
		AuthnProvider: userAuth,
	}

	ctx = context.WithValue(ctx, identity.IdentityService, service) // context passed into client rest methods

	client, clientErr := secretshub.NewClientWithResponses(*shurl,
		secretshub.WithRequestEditorFn(userAuth.Intercept))
	if clientErr != nil {
		slog.Error(fmt.Sprintf("failed to create secretshub client: %v", clientErr))
		os.Exit(1)
	}

	if client == nil {
		slog.Error("failed to create secretshub client")
		os.Exit(1)
	}

	slog.Info("Successfully created Secrets Hub client.")

	// Use the client to interact with the Secrets Hub API
}
