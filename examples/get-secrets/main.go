package main

import (
	"context"
	"encoding/json"
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

	filter := flag.String("filter", "", "Filter to apply when retrieving secrets")
	limit := flag.Int("limit", 100, "Maximum number of secrets to retrieve (default: 100, min: 1, max: 1000)")
	offset := flag.Int("offset", 0, "Offset for pagination (default: 0)")
	projectionFlag := flag.Bool("x", false, "Set Projection to get extended output")

	getall := flag.Bool("a", false, "Get all secrets, ignoring limit and offset")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if *limit > 1000 || *limit < 1 {
		slog.Error("Limit must be between 1 and 1000")
		os.Exit(1)
	}

	slogOpts := &slog.HandlerOptions{}
	if *debug {
		slogOpts.Level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, slogOpts).WithAttrs([]slog.Attr{
		slog.String("service", "secretshub"),
		slog.String("command", "sh"),
		slog.String("subcommand", "get-secrets"),
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

	ctx = context.WithValue(ctx, identity.IdentityService, service)

	client, clientErr := secretshub.NewClientWithResponses(*shurl, secretshub.WithRequestEditorFn(userAuth.Intercept))
	if clientErr != nil {
		slog.Error(fmt.Sprintf("failed to create secretshub client: %v", clientErr))
		os.Exit(1)
	}

	// curl --request GET \
	//  --url https://virtserver.swaggerhub.com/CyberArk/secretshub-api/fed-id/api/secrets \
	//  --header 'Accept: application/json' \
	//  --header 'Accept: application/x.secretshub.beta+json' \
	//  --header 'Authorization: Bearer 123'

	var secretList *secretshub.SecretListOutput
	var err error

	projection := "REGULAR"
	if *projectionFlag {
		projection = "EXTEND"
	}

	if *getall {
		secretList, err = secretshub.FetchAllSecrets(ctx, client, projection, *filter)
	} else {
		secretList, err = secretshub.GetSecretsPage(ctx, client, projection, *filter, *limit, *offset)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("failed to fetch secrets: %v", err))
		os.Exit(1)
	}

	PrintAsJSON(secretList.Secrets)

	slog.Debug("Successfully listed secret store secrets")
}

func PrintAsJSON(obj interface{}) {
	if obj == nil {
		return
	}
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal object to JSON: %v\n", err)
		return
	}
	if string(data) == "null" {
		return
	}
	fmt.Println(string(data))
}
