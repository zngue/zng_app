package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/zngue/zng_app/core/errors_ez"
	"github.com/zngue/zng_app/core/http/core"
	"github.com/zngue/zng_app/core/http/server"
)

type Client interface {
	Invoke(ctx context.Context, route ClientRoute, req any, reply any, opts ...CallOption) error
}

type Resolver func(ctx context.Context, serviceName string) (string, error)

type ClientRoute struct {
	Operation string
	Method    string
	Path      string
	Body      string
}

type CallInfo struct {
	BaseURL string
	Headers map[string]string
	Query   url.Values
}

type CallOption func(*CallInfo)

func WithBaseURL(baseURL string) CallOption {
	return func(info *CallInfo) {
		info.BaseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHeader(key string, value string) CallOption {
	return func(info *CallInfo) {
		if info.Headers == nil {
			info.Headers = make(map[string]string)
		}
		info.Headers[key] = value
	}
}

func WithQuery(key string, value string) CallOption {
	return func(info *CallInfo) {
		if info.Query == nil {
			info.Query = make(url.Values)
		}
		info.Query.Set(key, value)
	}
}

func InvokeUnary[Req any, Reply any](ctx context.Context, client Client, route ClientRoute, req *Req, opts ...CallOption) (*Reply, error) {
	reply := new(Reply)
	if err := client.Invoke(ctx, route, req, reply, opts...); err != nil {
		return nil, err
	}
	return reply, nil
}

func InvokeStream[Event any](ctx context.Context, client *JSONClient, route ClientRoute, req any, opts []CallOption, onEvent func(Event) error) error {
	baseURL, err := client.resolveBaseURL(ctx)
	if err != nil {
		return err
	}
	call := CallInfo{
		BaseURL: baseURL,
		Headers: make(map[string]string, len(client.headers)),
		Query:   make(url.Values),
	}
	for key, value := range client.headers {
		call.Headers[key] = value
	}
	for _, option := range opts {
		option(&call)
	}

	requestID := strings.TrimSpace(core.RequestIDFromContext(ctx))
	if requestID == "" {
		requestID = server.NewRequestID()
		ctx = core.WithRequestID(ctx, requestID)
	}
	if _, exists := call.Headers[core.RequestIDHeader]; !exists {
		call.Headers[core.RequestIDHeader] = requestID
	}

	requestURL, err := url.JoinPath(call.BaseURL, route.Path)
	if err != nil {
		return errors_ez.WrapF(err, "build request url")
	}
	encoded, err := core.EncodeQuery(req)
	if err != nil {
		return errors_ez.WrapF(err, "encode query")
	}
	for key, items := range encoded {
		if _, exists := call.Query[key]; !exists {
			for _, item := range items {
				call.Query.Add(key, item)
			}
		}
	}
	if len(call.Query) > 0 {
		requestURL += "?" + call.Query.Encode()
	}

	httpReq, err := http.NewRequestWithContext(ctx, route.Method, requestURL, nil)
	if err != nil {
		return errors_ez.WrapF(err, "create request")
	}
	for key, value := range call.Headers {
		httpReq.Header.Set(key, value)
	}
	httpReq.Header.Set("X-Operation", route.Operation)

	httpResp, err := client.httpClient.Do(httpReq)
	if err != nil {
		return errors_ez.WrapF(err, "do request")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return errors_ez.New("http status %d", httpResp.StatusCode)
	}

	scanner := bufio.NewScanner(httpResp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var event Event
		if err := json.Unmarshal(line, &event); err != nil {
			return errors_ez.WrapF(err, "decode stream event")
		}
		if err := onEvent(event); err != nil {
			return err
		}
	}
	return scanner.Err()
}

type JSONClient struct {
	baseURL     string
	resolver    Resolver
	serviceName string
	httpClient  *http.Client
	headers     map[string]string
}

type JSONClientOption func(*JSONClient)

func NewJSONClient(baseURL string, options ...JSONClientOption) *JSONClient {
	client := &JSONClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: http.DefaultClient,
		headers: map[string]string{
			"Accept": "application/json",
		},
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func NewJSONClientWithResolver(resolver Resolver, serviceName string, options ...JSONClientOption) *JSONClient {
	client := &JSONClient{
		resolver:    resolver,
		serviceName: serviceName,
		httpClient:  http.DefaultClient,
		headers: map[string]string{
			"Accept": "application/json",
		},
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (c *JSONClient) resolveBaseURL(ctx context.Context) (string, error) {
	if c.resolver != nil {
		addr, err := c.resolver(ctx, c.serviceName)
		if err != nil {
			return "", errors_ez.WrapF(err, "resolve service %s", c.serviceName)
		}
		return "http://" + strings.TrimRight(addr, "/"), nil
	}
	return c.baseURL, nil
}

func WithHTTPClient(httpClient *http.Client) JSONClientOption {
	return func(client *JSONClient) {
		client.httpClient = httpClient
	}
}

func WithDefaultHeader(key string, value string) JSONClientOption {
	return func(client *JSONClient) {
		client.headers[key] = value
	}
}

func (c *JSONClient) Invoke(ctx context.Context, route ClientRoute, req any, reply any, opts ...CallOption) error {
	baseURL, err := c.resolveBaseURL(ctx)
	if err != nil {
		return err
	}
	call := CallInfo{
		BaseURL: baseURL,
		Headers: make(map[string]string, len(c.headers)),
		Query:   make(url.Values),
	}
	for key, value := range c.headers {
		call.Headers[key] = value
	}
	for _, option := range opts {
		option(&call)
	}

	requestID := strings.TrimSpace(core.RequestIDFromContext(ctx))
	if requestID == "" {
		requestID = server.NewRequestID()
		ctx = core.WithRequestID(ctx, requestID)
	}
	if _, exists := call.Headers[core.RequestIDHeader]; !exists {
		call.Headers[core.RequestIDHeader] = requestID
	}

	var body io.Reader
	if route.Body != "" {
		payload, err := json.Marshal(req)
		if err != nil {
			return errors_ez.WrapF(err, "marshal request body")
		}
		body = bytes.NewReader(payload)
		call.Headers["Content-Type"] = "application/json"
	} else {
		encoded, err := core.EncodeQuery(req)
		if err != nil {
			return errors_ez.WrapF(err, "encode query")
		}
		for key, items := range encoded {
			if _, exists := call.Query[key]; !exists {
				for _, item := range items {
					call.Query.Add(key, item)
				}
			}
		}
	}

	requestURL, err := url.JoinPath(call.BaseURL, route.Path)
	if err != nil {
		return errors_ez.WrapF(err, "build request url")
	}
	if len(call.Query) > 0 {
		requestURL += "?" + call.Query.Encode()
	}

	httpReq, err := http.NewRequestWithContext(ctx, route.Method, requestURL, body)
	if err != nil {
		return errors_ez.WrapF(err, "create request")
	}
	for key, value := range call.Headers {
		httpReq.Header.Set(key, value)
	}
	httpReq.Header.Set("X-Operation", route.Operation)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return errors_ez.WrapF(err, "do request")
	}
	defer httpResp.Body.Close()

	var envelope core.Response
	envelope.Data = reply
	if err := json.NewDecoder(httpResp.Body).Decode(&envelope); err != nil {
		return errors_ez.WrapF(err, "decode response")
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		stack := append(errors_ez.CaptureStack(2), envelope.Stack...)
		return &errors_ez.Error{
			Code:    envelope.StatusCode,
			Reason:  envelope.Reason,
			Message: envelope.Message,
			Stack:   stack,
		}
	}
	return nil
}
