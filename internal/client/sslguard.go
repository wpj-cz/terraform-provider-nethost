package client

import "context"

type SSLGuardItem struct {
	ID     string
	Domain string
}

func (c *Client) CreateSSLGuard(ctx context.Context, domain string, port int64) error {
	payload := map[string]any{
		"domain": domain,
		"port":   port,
	}

	return c.post(ctx, "/services/sslguard/create", payload, nil)
}

func (c *Client) FindSSLGuardByDomain(ctx context.Context, domain string) (*SSLGuardItem, error) {
	var rawItems []struct {
		ID     any    `json:"id"`
		Domain string `json:"domain"`
	}

	err := c.post(ctx, "/services/sslguard/search", nil, &rawItems)
	if err != nil {
		return nil, err
	}

	for _, item := range rawItems {
		if item.Domain != domain {
			continue
		}

		return &SSLGuardItem{
			ID:     stringifyID(item.ID),
			Domain: item.Domain,
		}, nil
	}

	return nil, nil
}

func (c *Client) DeleteSSLGuard(ctx context.Context, id string) error {
	payload := map[string]any{
		"id": id,
	}

	return c.post(ctx, "/services/sslguard/delete", payload, nil)
}
