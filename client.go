package go_huawei

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/karmadon/go-huawei/internal"
	"github.com/karmadon/go-huawei/metrics"
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

const (
	ExperienceIdHeaderName = "X-GOOG-MAPS-EXPERIENCE-ID"
)

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

func (c *Client) get(ctx context.Context, config *apiConfig, apiReq interface{}, routeService RouteService) (*http.Response, error) {
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

func (c *Client) post(ctx context.Context, config *apiConfig, apiReq interface{}, routeService RouteService) (*http.Response, error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", host+config.path+string(routeService), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	q, err := c.generateAuthQuery(config.path, config.acceptsApiKey)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = q
	return c.do(ctx, req)
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := c.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	return client.Do(req.WithContext(ctx))
}

func (c *Client) getJSON(ctx context.Context, config *apiConfig, apiReq *DirectionsRequest, resp interface{}, routeService RouteService) error {
	requestMetrics := c.metricReporter.NewRequest(config.path)
	httpResp, err := c.get(ctx, config, apiReq, routeService)
	if err != nil {
		requestMetrics.EndRequest(ctx, err, httpResp, "")
		return err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(resp)
	requestMetrics.EndRequest(ctx, err, httpResp, httpResp.Header.Get("Server"))
	return err
}

func (c *Client) postJSON(ctx context.Context, config *apiConfig, apiReq interface{}, resp interface{}, routeService RouteService) error {
	requestMetrics := c.metricReporter.NewRequest(config.path)
	httpResp, err := c.post(ctx, config, apiReq, routeService)
	if err != nil {
		requestMetrics.EndRequest(ctx, err, httpResp, "")
		return err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(resp)
	requestMetrics.EndRequest(ctx, err, httpResp, httpResp.Header.Get("x-goog-maps-metro-area"))
	return err
}

type binaryResponse struct {
	statusCode  int
	contentType string
	data        io.ReadCloser
}

func (c *Client) generateAuthQuery(path string, acceptsApiKey bool) (string, error) {
	if c.apiKey != "" {
		if acceptsApiKey && len(c.apiKey) > 0 {
			return internal.SignURLWithApiKey(c.apiKey)
		}
	}

	return "", errors.New("map-kit: API Key missing")
}

// commonResponse contains the common response fields to most API calls
type commonResponse struct {
	// Status contains the status of the request, and may contain debugging
	// information to help you track down why the call failed.
	Status string `json:"status"`

	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// StatusError returns an error if this object has a Status different
// from OK or ZERO_RESULTS.
func (c *commonResponse) StatusError() error {
	if c.Status != "OK" && c.Status != "ZERO_RESULTS" {
		return fmt.Errorf("map-kit: %s - %s", c.Status, c.ErrorMessage)
	}
	return nil
}
