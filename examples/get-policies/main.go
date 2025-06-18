package main

import (
	"context"
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

	safename := flag.String("safename", "", "Safe name to filter policies by")

	ver := flag.Bool("version", false, "Print version")
	debug := flag.Bool("d", false, "Enable debug settings")
	flag.Parse()

	if *ver {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if *idtenanturl == "" || *iduser == "" || *idpass == "" || *shurl == "" || *safename == "" {
		flag.Usage()
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

	//curl --request GET \
	//--url 'https://{host}/api/policies?filter=filter.safeName+EQ+MySafeName&projection=REGULAR' \
	//--header 'Accept: application/json' \
	//--header 'Authorization: Bearer 123'

	projection := "EXTEND" //  Can be EXTEND, REGULAR, METADATA.
	filter := fmt.Sprintf("filter.safeName EQ %s", *safename)

	//safeName	Filter the sync policies by the Safe name.
	//GET https://<sub domain>.secretshub.cyberark.cloud/api/policies?projection=EXTEND&filter=filter.safeName EQ MySafeName
	//
	//target.id	Filter the sync policies that are syncing to a specific target secret store by its Secrets Hub ID.
	//GET https://<sub domain>.secretshub.cyberark.cloud/api/policies?filter=target.id EQ store-cfd25162-f8a9-4d94-8d36-f46c4b60d651

	params := &secretshub.ListPoliciesApiPoliciesGetParams{
		Projection: projection,
		Filter:     filter,
	}
	policies, err := client.ListPoliciesApiPoliciesGetWithResponse(ctx, params, AddAcceptApplicationJSONHeader)
	if err != nil {
		logger.Fatalf("failed to get policies: %v", err)
	}
	if policies.JSON400 != nil {
		logger.Fatalf("failed to get policies, error: %s", policies.JSON400.String())
	}
	if policies.JSON401 != nil {
		logger.Fatalf("failed to get policies, error: %s", policies.JSON401.String())
	}
	if policies.JSON403 != nil {
		logger.Fatalf("failed to get policies, error: %s", policies.JSON403.String())
	}
	if policies.JSON500 != nil {
		logger.Fatalf("failed to get policies, error: %s", policies.JSON500.String())
	}
	for pp := range policies.JSON200.Policies {
		p, err := policies.JSON200.Policies[pp].AsPolicyExtendedOutput()
		if err != nil {
			logger.Fatalf("failed to convert policy: %v", err)
		}
		fmt.Printf("Policy: %+v\n", p)
	}

	logger.Println("Successfully listed sync policies")
}

func AddAcceptApplicationJSONHeader(_ context.Context, req *http.Request) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/x.secretshub.beta+json")
	return nil
}
