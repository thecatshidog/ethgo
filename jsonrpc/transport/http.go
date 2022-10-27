package transport

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"time"

	"github.com/thecatshidog/ethgo/jsonrpc/codec"
	"github.com/valyala/fasthttp"
)

// HTTP is an http transport
type HTTP struct {
	addr    string
	client  *fasthttp.Client
	headers map[string]string
}

func newHTTP(addr string, headers map[string]string) *HTTP {
	return &HTTP{
		addr: addr,
		client: &fasthttp.Client{
			MaxConnsPerHost:     50000,
			ReadTimeout:         time.Second * 1,
			MaxIdleConnDuration: time.Second * 1000 * 60,
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(addr, time.Second*2)
			},
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
		},
		headers: headers,
	}
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	// Encode json-rpc request
	request := codec.Request{
		JsonRPC: "2.0",
		Method:  method,
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request.Params = data
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.addr)

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	req.SetBody(raw)

	if err := h.client.DoTimeout(req, res, 1*time.Second); err != nil {
		return err
	}

	// Decode json-rpc response
	var response codec.Response
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}

// SetMaxConnsPerHost sets the maximum number of connections that can be established with a host
func (h *HTTP) SetMaxConnsPerHost(count int) {
	h.client.MaxConnsPerHost = count
}
