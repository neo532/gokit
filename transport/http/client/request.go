package client

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/neo532/gokit/crypt/marshaler"
	"github.com/neo532/gokit/metadata"
	"github.com/neo532/gokit/middleware"
	xhttp "github.com/neo532/gokit/transport/http"
)

const (
	HeaderKeyStatusCode = "StatusCode"
	HeaderKeyUrl        = "Url"
	HeaderKeyCost       = "Cost"
	HeaderKeyLimit      = "Limit"
)

type Request struct {
	clt Client

	url                 string
	method              string
	contentType         string
	contentTypeResponse string

	retryTimes       int
	retryMaxDuration time.Duration
	retryDuration    time.Duration

	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
}

// ========== Opt ==========
type RequestOption func(o *Request)

func WithTimeLimit(d time.Duration) RequestOption {
	return func(o *Request) {
		o.clt = o.clt.Copy()
		o.clt.HttpClient().Timeout = d
	}
}
func WithUrl(s string) RequestOption {
	return func(o *Request) {
		o.url = s
	}
}
func WithMethod(m string) RequestOption {
	return func(o *Request) {
		o.method = m
	}
}
func WithContentType(ct string) RequestOption {
	return func(o *Request) {
		o.contentType = ct
	}
}
func WithContentTypeResponse(ct string) RequestOption {
	return func(o *Request) {
		o.contentTypeResponse = ct
	}
}
func WithRetryTimes(times int) RequestOption {
	return func(o *Request) {
		o.retryTimes = times
	}
}
func WithRetryDuration(d time.Duration) RequestOption {
	return func(o *Request) {
		o.retryDuration = d
	}
}
func WithRetryMaxDuration(d time.Duration) RequestOption {
	return func(o *Request) {
		o.retryMaxDuration = d
	}
}
func WithRequestEncoder(encoder EncodeRequestFunc) RequestOption {
	return func(o *Request) {
		o.encoder = encoder
	}
}
func WithResponseDecoder(decoder DecodeResponseFunc) RequestOption {
	return func(o *Request) {
		o.decoder = decoder
	}
}
func WithErrorDecoder(errorDecoder DecodeErrorFunc) RequestOption {
	return func(o *Request) {
		o.errorDecoder = errorDecoder
	}
}

// ========== /Opt ==========

func NewRequest(clt Client, opts ...RequestOption) (req *Request) {
	req = &Request{
		retryTimes:       clt.RetryTime(),
		retryDuration:    time.Microsecond,
		retryMaxDuration: 20 * time.Microsecond,

		errorDecoder: DefaultErrorDecoder,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		clt:          clt,
	}
	for _, o := range opts {
		o(req)
	}
	return
}

func (r *Request) Do(ctx context.Context, req interface{}, reply interface{}) (c context.Context, err error) {
	c = ctx

	h := func(ctx context.Context, req interface{}, reply interface{}) (c context.Context, err error) {
		c = ctx

		url := r.url
		if qa, ok := xhttp.FromClientContext(c); ok {
			url = qa.AppendToUrl(r.url)
		}

		var reqBody []byte
		if reqBody, err = r.encoder(c, r.contentType, req); err != nil {
			return
		}

		reqHeader, headerBCurl := r.FmtHeader(c)

		retryDuration := r.retryDuration
		var er error
		for i := 0; i <= r.retryTimes; i++ {

			var param *http.Request
			if param, err = http.NewRequestWithContext(
				c,
				r.method,
				url,
				bytes.NewReader(reqBody)); err != nil {
				return
			}
			param.Header = reqHeader

			// request
			var resp *http.Response
			start := time.Now()
			resp, err = r.clt.HttpClient().Do(param)
			cost := time.Now().Sub(start)

			var respCode int
			var respBody []byte
			var cancelRetry bool
			for j := 0; j < 1; j++ {
				if err != nil {
					break
				}
				if resp != nil {
					respCode = resp.StatusCode
					if resp.Body != nil {
						defer resp.Body.Close()
					}
				}
				if cancelRetry, err = r.errorDecoder(c, resp); err != nil {
					break
				}
				if resp != nil {
					if r.contentTypeResponse != "" {
						resp.Header.Set(ContentTypeHeaderKey, r.contentTypeResponse)
					}
					respBody, er = r.decoder(c, resp, reply)
				}
			}

			r.log(c, url, headerBCurl, reqBody, respCode, respBody, cost, err)

			md := metadata.New(nil)
			if resp != nil {
				md = metadata.New(resp.Header)
			}
			md.Set(HeaderKeyStatusCode, strconv.Itoa(respCode))
			md.Set(HeaderKeyUrl, r.url)
			md.Set(HeaderKeyCost, cost.String())
			md.Set(HeaderKeyLimit, r.clt.HttpClient().Timeout.String())
			c = metadata.NewClientResponseContext(c, md)

			if cancelRetry || err == nil {
				break
			}

			time.Sleep(r.retryDuration)
			if retryDuration < r.retryMaxDuration {
				retryDuration = retryDuration + retryDuration
			}
		}
		if err == nil {
			err = er
		}
		return
	}

	if len(r.clt.Middlewares()) > 0 {
		h = middleware.Chain(r.clt.Middlewares()...)(h)
	}
	return h(c, req, reply)
}

func (r *Request) log(c context.Context,
	url string, header strings.Builder, reqBody []byte,
	respCode int, respBody []byte, cost time.Duration, err error) {
	respStr := string(respBody)
	if l := r.clt.ResponseMaxLength(); l > 0 && utf8.RuneCountInString(respStr) > l {
		respStr = string([]rune(respStr)[:l]) + "..."
	}

	reqBodyS := string(reqBody)
	if len(reqBodyS) != 0 && reqBodyS != "{}" {
		reqBodyS = " -d '" + string(reqBodyS) + "'"
	} else {
		reqBodyS = ""
	}
	msg := fmt.Sprintf("[code:%d] [limit:%s] [cost:%s] [curl -X '%s' '%s'%s%s%s] [rst:%s]",
		respCode,
		r.clt.HttpClient().Timeout.String(),
		cost.String(),
		r.method,
		url,
		header.String(),
		r.clt.CurlArgs(),
		reqBodyS,
		respStr,
	)
	if err != nil {
		r.clt.Logger().Error(c, fmt.Sprintf("[err:%s] %s", err.Error(), msg))
		return
	}
	r.clt.Logger().Debug(c, msg)
}

func (r *Request) FmtHeader(c context.Context) (h http.Header, curl strings.Builder) {
	h = http.Header{}
	if md, ok := metadata.FromClientContext(c); ok {
		md.Range(func(k string, vs []string) (b bool) {
			for _, v := range vs {
				h.Set(k, v)
				curl.WriteString(" -H '" + k + ":" + v + "'")
			}
			return true
		})
	}
	if r.contentType != "" {
		h.Set(ContentTypeHeaderKey, r.contentType)
		curl.WriteString(" -H '" + ContentTypeHeaderKey + ":" + r.contentType + "'")
		return
	}
	if h.Get(ContentTypeHeaderKey) == "" && HasBody(r.method) {
		h.Set(ContentTypeHeaderKey, ContentTypeHeaderDefaultValue)
		curl.WriteString(" -H '" + ContentTypeHeaderKey + ":" + ContentTypeHeaderDefaultValue + "'")
	}
	return
}

func HasBody(method string) (b bool) {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		return
	case http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	}
	return true
}

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(c context.Context, res *http.Response) (cancelRetry bool, err error)

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(c context.Context, contentType string, in interface{}) (body []byte, err error)

// DecodeResponseFunc is response decode func.
type DecodeResponseFunc func(c context.Context, res *http.Response, out interface{}) (body []byte, err error)

// DefaultRequestEncoder is an HTTP request encoder.
func DefaultRequestEncoder(c context.Context, contentType string, in interface{}) (body []byte, err error) {
	subContentType := ContentSubtype(contentType)
	codec := marshaler.GetMarshaler(subContentType)
	if codec == nil {
		err = fmt.Errorf("Wrong content-type(%s) from header", subContentType)
		return
	}
	return codec.Marshal(in)
}

// DefaultResponseDecoder is an HTTP response decoder.
func DefaultResponseDecoder(c context.Context, res *http.Response, v interface{}) (body []byte, err error) {
	if v == nil {
		return
	}
	subContentType := ContentSubtype(res.Header.Get("Content-Type"))
	codec := marshaler.GetMarshaler(subContentType)
	if codec == nil {
		err = fmt.Errorf("Wrong content-type(%s) from header", subContentType)
		return
	}

	if body, err = io.ReadAll(res.Body); err != nil {
		return
	}
	err = codec.Unmarshal(body, v)
	return
}

// DefaultErrorDecoder is an HTTP error decoder.
func DefaultErrorDecoder(c context.Context, resp *http.Response) (cancelRetry bool, err error) {
	if resp == nil {
		err = errors.New("nil *http.Response")
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}
	if resp.StatusCode >= 400 && resp.StatusCode <= 407 {
		cancelRetry = true
		err = errors.New(resp.Status)
		return
	}
	err = errors.New(resp.Status)
	return
}
