package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal"
)

type httpClient struct {
	c       *http.Client
	baseURL string
	header  http.Header
}

func newHTTP(opt Options) *httpClient {
	baseURL := fmt.Sprintf(
		"%s://%s:%d/%s/",
		opt.Scheme, opt.HTTPHost, opt.HTTPPort, opt.Version,
	)
	return &httpClient{
		c:       &http.Client{},
		baseURL: baseURL,
		header:  opt.Header,
	}
}

func (c *httpClient) do(ctx context.Context, req internal.Endpoint, dest any) error {
	var body io.Reader
	if b := req.Body(); b != nil {
		marshaled, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(marshaled)
	}

	url := c.url(req)
	httpreq, err := http.NewRequestWithContext(ctx, req.Method(), url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Clone default request headers.
	httpreq.Header = c.header.Clone()

	res, err := c.c.Do(httpreq)
	if err != nil {
		return err
	}

	// Response body SHOULD always be read completely and closed
	// to allow the underlying [http.Transport] to re-use the TCP connection.
	// See: https://pkg.go.dev/net/http#Client.Do
	resBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	// TODO(dyma): not sure if we should always report this error.
	// What if we don't need the body because dest=nil and status is OK?
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode > 299 {
		if res.StatusCode == http.StatusNotFound {
			return nil // leave dest a nil
		}
		// TODO(dyma): better error handling?
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, resBody)
	}

	if dest != nil {
		if err := json.Unmarshal(resBody, dest); err != nil {
			return fmt.Errorf("unmarshal response body: %w", err)
		}
	}
	return nil
}

func (c *httpClient) url(req internal.Endpoint) string {
	var url strings.Builder

	url.WriteString(c.baseURL)
	url.WriteString(strings.TrimLeft(req.Path(), "/"))

	if query := req.Query(); len(query) > 0 {
		url.WriteString("?")
		url.WriteString(query.Encode())
	}

	return url.String()
}
