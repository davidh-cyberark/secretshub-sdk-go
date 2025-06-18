package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
	"github.com/davidh-cyberark/secretshub-sdk-go/secretshub"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	version = "dev"
)

func main() {
	idtenanturl := flag.String("idtenanturl", "", "Identity URL")
	iduser := flag.String("iduser", "", "Identity user id")
	idpass := flag.String("idpass", "", "Identity user password")
	shurl := flag.String("shurl", "", "Secrets Hub URL, Ex: https://EXAMPLE.secretshub.cyberark.cloud/")

	storepath := flag.String("storepath", "stores", "Path to secret stores data")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if storepath == nil || *storepath == "" {
		fmt.Println("Please provide a path to the secret stores data using -storepath flag")
		os.Exit(1)
	}

	// logger
	logger := log.New(os.Stderr, "[get-secrets] ", log.LstdFlags)

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

	client, clientErr := secretshub.NewClientWithResponses(*shurl, secretshub.WithRequestEditorFn(userAuth.Intercept))
	if clientErr != nil {
		logger.Fatalf("failed to create secretshub client: %v", clientErr)
	}

	// curl --request GET \
	//  --url https://virtserver.swaggerhub.com/CyberArk/secretshub-api/fed-id/api/secrets \
	//  --header 'Accept: application/json' \
	//  --header 'Accept: application/x.secretshub.beta+json' \
	//  --header 'Authorization: Bearer 123'
	store := &secretshub.SecretStoreWithReplicatedDataOutput{}
	err := readJSONFile(*storepath, store)
	if err != nil {
		logger.Fatalf("failed to read secret store from file: %v", err)
	}
	projection := "EXTEND"
	limit := 1000
	filterByStoreId := fmt.Sprintf("storeId CONTAINS %s", store.Id)
	params := &secretshub.ListSecretsApiSecretsGetParams{
		Projection: &projection,
		Limit:      &limit,
		Filter:     &filterByStoreId,
	}
	resp, err := client.ListSecretsApiSecretsGetWithResponse(ctx, params, AddAcceptApplicationJSONHeader)
	if err != nil {
		logger.Fatalf("failed to list secrets: %v", err)
	}
	for ss := range resp.JSON200.Secrets {
		s, err := resp.JSON200.Secrets[ss].AsExtendedSecretOutput()
		if err != nil {
			fmt.Printf("failed to convert secret to extended output: %v", err)
			continue
		}
		fmt.Printf("%+v\n", s)
	}
	logger.Println("Successfully listed secret store secrets")
}

func AddAcceptApplicationJSONHeader(ctx context.Context, req *http.Request) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/x.secretshub.beta+json")
	return nil
}

func readJSONFile(filePath string, target interface{}) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v", err)
		}
	}()

	// Read all the file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the json into the target variable
	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}
