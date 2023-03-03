package go_huawei

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/stremovskyy/go-huawei/internal"
	"github.com/stremovskyy/go-huawei/metrics"
)

type Client struct {
	httpClient        *http.Client
	apiKey            string
	baseURL           string
	requestsPerSecond int
	rateLimiter       *rate.Limiter
	metricReporter    metrics.Reporter
}

// ClientOption is the type of constructor options for NewClient(...).
type ClientOption func(*Client) error

var defaultRequestsPerSecond = 50

func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{
		requestsPerSecond: defaultRequestsPerSecond,
		metricReporter:    metrics.NoOpReporter{},
	}

	err := WithHTTPClient(&http.Client{})(c)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	if c.apiKey == "" {
		return nil, errors.New("map-kit: API Key or Maps for Work credentials missing")
	}

	if c.requestsPerSecond > 0 {
		c.rateLimiter = rate.NewLimiter(rate.Limit(c.requestsPerSecond), c.requestsPerSecond)
	}

	return c, nil
}

// WithHTTPClient configures a Maps API client with a http.Client to make requests
// over.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		if _, ok := c.Transport.(*transport); !ok {
			t := c.Transport
			if t != nil {
				c.Transport = &transport{Base: t}
			} else {
				//proxyUrl, _ := url.Parse("http://localhost:8888")
				//c.Transport = &transport{Base: &http.Transport{
				//	Proxy: http.ProxyURL(proxyUrl),
				//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				//},}

				c.Transport = http.DefaultTransport
			}
		}
		client.httpClient = c
		return nil
	}
}

// WithAPIKey configures a Maps API client with an API Key
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

// WithBaseURL configures a Maps API client with a custom base url
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// WithRateLimit configures the rate limit for back end requests. Default is to
// limit to 50 requests per second. A value of zero disables rate limiting.
func WithRateLimit(requestsPerSecond int) ClientOption {
	return func(c *Client) error {
		c.requestsPerSecond = requestsPerSecond
		return nil
	}
}

func WithMetricReporter(reporter metrics.Reporter) ClientOption {
	return func(c *Client) error {
		c.metricReporter = reporter
		return nil
	}
}

type apiConfig struct {
	host          string
	path          string
	acceptsApiKey bool
}

func (c *Client) awaitRateLimiter(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	return c.rateLimiter.Wait(ctx)
}

func (c *Client) get(ctx context.Context, config *apiConfig, _ interface{}, _ RouteService) ([]byte, error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	req, err := http.NewRequest("GET", host+config.path, nil)
	if err != nil {
		return nil, err
	}

	q, err := c.generateAuthQuery(config.path, config.acceptsApiKey)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q
	return c.do(ctx, req)
}

func (c *Client) post(ctx context.Context, config *apiConfig, apiReq interface{}, routeService RouteService) ([]byte, error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, NewGoHuaweiError("marshal json", err)
	}

	req, err := http.NewRequest("POST", host+config.path+string(routeService), bytes.NewBuffer(body))
	if err != nil {
		e := NewGoHuaweiError("post request", err)
		e.AddRawRequest(body)

		return nil, e
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	q, err := c.generateAuthQuery(config.path, config.acceptsApiKey)
	if err != nil {
		e := NewGoHuaweiError("post request", err)
		e.AddRawRequest(body)
		return nil, err
	}

	req.URL.RawQuery = q
	return c.do(ctx, req)
}

func (c *Client) do(ctx context.Context, req *http.Request) ([]byte, error) {
	client := c.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}

		defer reader.Close()
	default:
		reader = resp.Body
	}

	return ioutil.ReadAll(reader)
}

func (c *Client) getJSON(ctx context.Context, config *apiConfig, apiReq *DirectionsRequest, resp interface{}, routeService RouteService) error {
	requestMetrics := c.metricReporter.NewRequest(config.path)

	httpResp, err := c.get(ctx, config, apiReq, routeService)
	if err != nil {
		requestMetrics.EndRequest(ctx, err, httpResp, "")
		return err
	}

	if httpResp == nil {
		e := NewGoHuaweiError("empty response", err)
		e.AddRawResponse(httpResp)
		return e
	}

	err = json.Unmarshal(httpResp, resp)
	if err != nil {
		e := NewGoHuaweiError("unmarshal response", err)
		e.AddRawResponse(httpResp)
		return err
	}

	requestMetrics.EndRequest(ctx, err, httpResp, "")

	return err
}

func (c *Client) postJSON(ctx context.Context, config *apiConfig, apiReq interface{}, resp interface{}, routeService RouteService) error {
	requestMetrics := c.metricReporter.NewRequest(config.path)

	httpResp, err := c.post(ctx, config, apiReq, routeService)
	if err != nil {
		requestMetrics.EndRequest(ctx, err, httpResp, "")
		return err
	}

	if httpResp == nil {
		e := NewGoHuaweiError("empty response", err)
		e.AddRawResponse(httpResp)
		return e
	}

	err = json.Unmarshal(httpResp, resp)
	if err != nil {
		e := NewGoHuaweiError("unmarshal response", err)
		e.AddRawResponse(httpResp)
		return err
	}

	requestMetrics.EndRequest(ctx, err, httpResp, "")

	return err
}

func (c *Client) generateAuthQuery(_ string, acceptsApiKey bool) (string, error) {
	if c.apiKey != "" {
		if acceptsApiKey && len(c.apiKey) > 0 {
			return internal.SignURLWithApiKey(c.apiKey)
		}
	}

	return "", errors.New("map-kit: API Key missing")
}
