package client

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/middleware"
)

type Client struct {
	logger                   logger.ILogger
	mapValue                 sync.Map
	middlewares              []middleware.Middleware
	defaultResponseMaxLength int
	defaultRetryTimes        int

	curlArgs string

	transport *http.Transport

	httpClient *http.Client
}

// ========== Opt ==========
type ClientOption func(o *Client)

// ---------- xhttp ----------
func WithLogger(l logger.ILogger) ClientOption {
	return func(o *Client) {
		o.logger = l
	}
}
func WithResponseMaxLength(l int) ClientOption {
	return func(o *Client) {
		o.defaultResponseMaxLength = l
	}
}
func WithMap(key, value string) ClientOption {
	return func(o *Client) {
		o.mapValue.Store(key, value)
	}
}

// ---------- request ----------
func WithMiddleware(ms ...middleware.Middleware) ClientOption {
	return func(o *Client) {
		o.middlewares = append(o.middlewares, ms...)
	}
}
func WithDefaultRetryTimes(times int) ClientOption {
	return func(o *Client) {
		o.defaultRetryTimes = times
	}
}
func WithDefaultTimeLimit(d time.Duration) ClientOption {
	return func(o *Client) {
		o.httpClient.Timeout = d
	}
}

// ---------- connect pool ----------
func WithIdleConnTimeout(d time.Duration) ClientOption {
	return func(o *Client) {
		o.transport.IdleConnTimeout = d
	}
}
func WithMaxConnsPerHost(n int) ClientOption {
	return func(o *Client) {
		o.transport.MaxConnsPerHost = n
	}
}
func WithMaxIdleConns(n int) ClientOption {
	return func(o *Client) {
		o.transport.MaxIdleConns = n
	}
}
func WithMaxIdleConnsPerHost(n int) ClientOption {
	return func(o *Client) {
		o.transport.MaxIdleConnsPerHost = n
	}
}

// ---------- tls ----------
func initTLS(o *Client) {
	if o.transport.TLSClientConfig == nil {
		o.transport.TLSClientConfig = &tls.Config{}
	}
}
func WithInsecureSkipVerify(b bool) ClientOption {
	return func(o *Client) {
		initTLS(o)
		o.transport.TLSClientConfig.InsecureSkipVerify = b
		if b {
			o.setCurlArgs(" --insecure")
		}
	}
}
func WithCertFile(crt, key string) (oR ClientOption, err error) {
	var cert tls.Certificate
	if cert, err = tls.LoadX509KeyPair(crt, key); err != nil {
		return
	}
	oR = func(o *Client) {
		initTLS(o)
		if o.transport.TLSClientConfig.Certificates == nil {
			o.transport.TLSClientConfig.Certificates = make([]tls.Certificate, 0, 1)
		}
		o.transport.TLSClientConfig.Certificates = append(
			o.transport.TLSClientConfig.Certificates,
			cert,
		)
		o.setCurlArgs(fmt.Sprintf(
			" --cert %s --key %s",
			strings.TrimSpace(crt),
			strings.TrimSpace(key),
		))
	}
	return
}
func WithCaCertFile(crt string) (oR ClientOption, err error) {
	var caCrt []byte
	if caCrt, err = os.ReadFile(crt); err != nil {
		return
	}

	oR = func(o *Client) {
		initTLS(o)
		if o.transport.TLSClientConfig.RootCAs == nil {
			o.transport.TLSClientConfig.RootCAs = x509.NewCertPool()
		}
		o.transport.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCrt)
		o.setCurlArgs(" --cacert " + strings.TrimSpace(crt))
	}
	return
}

// ========== /Opt ==========

func NewClient(opts ...ClientOption) (client Client) {

	client = Client{
		logger:            &logger.DefaultILogger{},
		middlewares:       make([]middleware.Middleware, 0, 10),
		defaultRetryTimes: 0,
		transport:         http.DefaultTransport.(*http.Transport).Clone(),
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	client.transport.MaxConnsPerHost = runtime.NumCPU()*2 + 1
	for _, o := range opts {
		o(&client)
	}
	client.httpClient.Transport = client.transport
	return
}

func (r *Client) setCurlArgs(arg string) {
	if strings.Contains(r.curlArgs, arg) == false {
		r.curlArgs += arg
	}
}

func (r *Client) CurlArgs() string {
	return r.curlArgs
}

func (r Client) Logger() logger.ILogger {
	return r.logger
}

func (r Client) ResponseMaxLength() int {
	return r.defaultResponseMaxLength
}

func (r Client) RetryTime() int {
	return r.defaultRetryTimes
}

func (r Client) Value(key string) (value string) {
	if v, ok := r.mapValue.Load(key); ok {
		if s, ok := v.(string); ok {
			value = s
		}
	}
	return
}

func (r Client) CopyMiddleware() Client {
	mds := make([]middleware.Middleware, 0, 1)
	for _, mw := range mds {
		mds = append(mds, mw)
	}
	r.middlewares = mds
	return r
}

func (r Client) Copy() (clt Client) {
	clt = Client{
		defaultResponseMaxLength: r.defaultResponseMaxLength,
		logger:                   r.logger,
		middlewares:              r.middlewares,
		defaultRetryTimes:        r.defaultRetryTimes,
		httpClient: &http.Client{
			Timeout:   3 * time.Second,
			Transport: r.httpClient.Transport,
		},
	}
	return
}

func (r Client) AddMiddleware(mds ...middleware.Middleware) Client {
	for _, mw := range mds {
		r.middlewares = append(r.middlewares, mw)
	}
	return r
}

func (r Client) Middlewares() (ms []middleware.Middleware) {
	return r.middlewares
}

func (r Client) HttpClient() *http.Client {
	return r.httpClient
}
