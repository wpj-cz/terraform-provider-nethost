package client

import (
	"context"
	"gitlab.wpj.cz/terraform-provider-nethost/internal/tfconvert"
)

type SpongefingerItem struct {
	ID        string
	Name      string
	Domain    string
	VariantID int64
}

type SpongefingerDetailItem struct {
	ID         string
	Name       string
	Domain     string
	VariantID  int64
	Subdomains []string
	Backends   []string
}

func (c *Client) FindSpongefingerByDomain(ctx context.Context, domain string) (*SpongefingerItem, error) {
	var rawItems []struct {
		ID      any    `json:"id"`
		Name    string `json:"name"`
		Domain  string `json:"domain"`
		Variant struct {
			ID any `json:"id"`
		} `json:"variant"`
	}

	err := c.post(ctx, "/services/spongefinger/search", nil, &rawItems)
	if err != nil {
		return nil, err
	}

	for _, item := range rawItems {
		if item.Domain != domain {
			continue
		}

		return &SpongefingerItem{
			ID:        stringifyID(item.ID),
			Name:      item.Name,
			Domain:    item.Domain,
			VariantID: tfconvert.Int64FromAny(item.Variant.ID),
		}, nil
	}

	return nil, nil
}

func (c *Client) CreateSpongefinger(ctx context.Context, name string, variantID int64, domain string, subdomains []string, backends []string) error {
	if subdomains == nil {
		subdomains = []string{}
	}
	if backends == nil {
		backends = []string{}
	}

	payload := map[string]any{
		"name":       name,
		"variant_id": variantID,
		"domain":     domain,
		"subdomains": subdomains,
		"backends":   backends,
	}

	return c.post(ctx, "/services/spongefinger/create", payload, nil)
}

func (c *Client) UpdateSpongefinger(ctx context.Context, id string, name string, domain string, subdomains []string, backends []string) error {
	if subdomains == nil {
		subdomains = []string{}
	}
	if backends == nil {
		backends = []string{}
	}

	payload := map[string]any{
		"id":         id,
		"name":       name,
		"domain":     domain,
		"subdomains": subdomains,
		"backends":   backends,
	}

	return c.post(ctx, "/services/spongefinger/update", payload, nil)
}

func (c *Client) FindSpongefingerByID(ctx context.Context, id string) (*SpongefingerDetailItem, error) {
	var rawItem struct {
		ID         any      `json:"id"`
		Name       string   `json:"name"`
		Domain     string   `json:"domain"`
		Subdomains []string `json:"subdomains"`
		Backends   []string `json:"backends"`
		Variant    struct {
			ID any `json:"id"`
		} `json:"variant"`
	}

	payload := map[string]any{
		"id": id,
	}

	err := c.post(ctx, "/services/spongefinger/detail", payload, &rawItem)
	if err != nil {
		return nil, err
	}

	return &SpongefingerDetailItem{
		ID:         stringifyID(rawItem.ID),
		Name:       rawItem.Name,
		Domain:     rawItem.Domain,
		VariantID:  tfconvert.Int64FromAny(rawItem.Variant.ID),
		Subdomains: rawItem.Subdomains,
		Backends:   rawItem.Backends,
	}, nil
}

func (c *Client) DeleteSpongefinger(ctx context.Context, id string) error {
	payload := map[string]any{
		"id": id,
	}

	return c.post(ctx, "/services/spongefinger/delete", payload, nil)
}
