package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	storedir := flag.String("storedir", "stores", "Directory to save secret stores data (default: stores)")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *idtenanturl == "" || *iduser == "" || *idpass == "" || *shurl == "" {
		fmt.Println("Please provide all required flags: -idtenanturl, -iduser, -idpass, -shurl")
		os.Exit(1)
	}

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if storedir == nil || *storedir == "" {
		fmt.Println("Please provide a directory to save secret stores data using -storedir flag")
		os.Exit(1)
	}
	curdir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	newpath := filepath.Join(curdir, *storedir)
	err = os.MkdirAll(newpath, 0755) // Creates the directory
	if err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	// logger
	logger := log.New(os.Stderr, "[get-all-stores] ", log.LstdFlags)

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

	var behavior secretshub.ListAllSecretStoresApiSecretStoresGetParamsBehavior = "SECRETS_TARGET"
	filter := "type EQ GCP_GSM"
	paramsGetAllSecretStores := &secretshub.ListAllSecretStoresApiSecretStoresGetParams{
		Behavior: behavior,
		Filter:   &filter,
	}
	storesResp, storesRespErr := client.ListAllSecretStoresApiSecretStoresGetWithResponse(ctx,
		paramsGetAllSecretStores,
		AddAcceptApplicationJSONHeader)

	if storesRespErr != nil {
		logger.Fatalf("failed to list secret stores: %v", storesRespErr)
	}
	//reg := regexp.MustCompile("[^a-zA-Z]") // Matches anything that is not a letter

	if storesResp.JSON200 != nil {
		for _, store := range storesResp.JSON200.SecretStores {
			SaveStoreDataToFile(*storedir, &store, logger)
		}
	}
	if storesResp.JSON400 != nil {
		logger.Fatalf("failed to list secret stores: %s", storesResp.JSON400.Message)
	}
	if storesResp.JSON401 != nil {
		logger.Fatalf("failed to list secret stores: %s", storesResp.JSON401.Message)
	}
	if storesResp.JSON403 != nil {
		logger.Fatalf("failed to list secret stores: %s", storesResp.JSON403.Message)
	}
	if storesResp.JSON500 != nil {
		logger.Fatalf("failed to list secret stores: %s", storesResp.JSON500.Message)
	}
	logger.Println("Successfully listed secret stores")
}

func AddAcceptApplicationJSONHeader(ctx context.Context, req *http.Request) error {
	req.Header.Add("Accept", "application/json")
	return nil
}

func SaveStoreDataToFile(dir string, store *secretshub.SecretStoreWithReplicatedDataOutput, logger *log.Logger) {
	if store == nil {
		logger.Println("Store is nil")
		return
	}

	storeFileName := fmt.Sprintf("%s.json", store.Id)
	fileName := filepath.Join(dir, storeFileName)
	file, err := os.Create(fileName)
	if err != nil {
		logger.Printf("failed to create file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		logger.Printf("failed to marshal data: %v", err)
		return
	}

	if _, err := file.Write(jsonData); err != nil {
		logger.Printf("failed to write data to file %s: %v", fileName, err)
		return
	}

	logger.Printf("Store data saved to %s\n", fileName)
}
