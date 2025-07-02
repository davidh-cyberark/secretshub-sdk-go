package secretshub

import (
	"context"
	"fmt"
	"math"
	"net/http"
)

// GetSecretsPage retrieves a page of secrets from the Secrets Hub API.
func GetSecretsPage(ctx context.Context, client *ClientWithResponses, projection string, filter string, limit int, offset int) (*SecretListOutput, error) {
	params := &ListSecretsApiSecretsGetParams{
		Limit:  &limit,
		Offset: &offset,
	}
	if len(projection) > 0 {
		params.Projection = &projection
	}
	if len(filter) > 0 {
		params.Filter = &filter
	}

	resp, err := client.ListSecretsApiSecretsGetWithResponse(ctx, params, AddAcceptApplicationJSONHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	if resp.JSON500 != nil {
		return nil, resp.JSON500
	}
	if resp.JSON400 != nil {
		return nil, resp.JSON400
	}
	if resp.JSON401 != nil {
		return nil, resp.JSON401
	}
	if resp.JSON403 != nil {
		return nil, resp.JSON403
	}
	if resp.JSON405 != nil {
		return nil, resp.JSON405
	}
	if resp.JSON406 != nil {
		return nil, resp.JSON406
	}

	if resp.ApplicationxSecretshubBetaJSON200 != nil {
		return resp.ApplicationxSecretshubBetaJSON200, nil
	}

	return resp.JSON200, nil
}

// FetchAllSecrets retrieves all secrets from the Secrets Hub API, handling pagination.
func FetchAllSecrets(ctx context.Context, client *ClientWithResponses, projection string, filter string) (*SecretListOutput, error) {
	var allSecrets []SecretListOutput_Secrets_Item

	count := 0
	total_count := math.MaxInt64
	for count < total_count {
		page, err := GetSecretsPage(ctx, client, projection, filter, 1000, count)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secrets: %w", err)
		}
		if page == nil {
			break
		}
		allSecrets = append(allSecrets, page.Secrets...)
		count += *page.Count
		total_count = *page.TotalCount
	}
	result := SecretListOutput{
		Secrets:    allSecrets,
		Count:      &count,
		TotalCount: &total_count,
	}

	return &result, nil
}

func AddAcceptApplicationJSONHeader(ctx context.Context, req *http.Request) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/x.secretshub.beta+json")
	return nil
}
