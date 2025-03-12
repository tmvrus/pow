package provider

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type Provider struct {
	host string
}

func NewProvider(host string) *Provider {
	return &Provider{
		host: host,
	}
}

func (p *Provider) Get(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.host, nil)
	if err != nil {
		return "", fmt.Errorf("create http request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do http request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", res.StatusCode)
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		return "", fmt.Errorf("parse html response: %w", err)
	}

	// I know, it looks ugly, but they do not provide any API.
	for n := range node.Descendants() {
		if n.Type == html.TextNode &&
			n.Parent != nil &&
			n.Parent.Type == html.ElementNode &&
			n.Parent.Data == "a" {

			return n.Data, nil
		}
	}

	return "", fmt.Errorf("invalid remote document format")
}
