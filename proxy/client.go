package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	c := &Client{
		client: &http.Client{
			Transport: &http2.Transport{
				// So http2.Transport doesn't complain the URL scheme isn't 'https'
				AllowHTTP: false,
				// Pretend we are dialing a TLS endpoint.
				// Note, we ignore the passed tls.Config
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		},
	}
	return c
}

func (c *Client) Post(host, path string, data []byte) ([]byte, error) {
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		},
		Header: http.Header{
			"content-type": []string{"application/grpc+json"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(c.PrepData(data))),
	}

	// Sends the request
	r, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, err
	}

	defer r.Body.Close()

	if r.Header.Get("Grpc-Status") != "" && r.Header.Get("Grpc-Status") != "0" {
		return nil, errors.New(r.Header.Get("Grpc-Message"))
	}

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return contents[4:], nil
}

func (c *Client) PrepData(data []byte) []byte {
	compAndLen := make([]byte, 5)
	// compress none
	compAndLen[0] = 0
	binary.BigEndian.PutUint32(compAndLen[1:], uint32(len(data)))

	return append(compAndLen, data...)
}
