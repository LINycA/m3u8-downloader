package request

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// Request 请求
type Request struct {
	Headers      map[string]string // 请求报头
	Method       string            // 请求方式
	Body         string            // 请求body
	Uri          string            // 请求地址
	ResponseBody []byte            // 请求回的数据
	RetryCount   int               // 请求重试次数
	Cli          *http.Client      // 可复用的client
}

// NewRequest 新建一个请求对象
func NewRequest(method, body, uri string, retryCount int, headers map[string]string, cli *http.Client) *Request {
	return &Request{
		Method:     method,
		Headers:    headers,
		Body:       body,
		Uri:        uri,
		Cli:        cli,
		RetryCount: retryCount,
	}
}

// DoRequest 发起请求，如存在错误则重试
func (r *Request) DoRequest() error {
	var err error
	for tryCount := 0; tryCount < r.RetryCount; tryCount++ {
		err = r.doRequest()
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		} else {
			return nil
		}
	}
	return err
}

// doRequest 发请求
func (r *Request) doRequest() error {
	body := bytes.NewReader([]byte(r.Body))
	req, err := http.NewRequest(r.Method, r.Uri, body)
	if err != nil {
		return err
	}
	if r.Headers != nil {
		if len(r.Headers) > 0 {
			for k, v := range r.Headers {
				req.Header.Add(k, v)
			}
		}
	}
	resp, err := r.Cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r.ResponseBody = respBody
	return nil
}

// NewClient 新建请求客户端
func NewClient(timeout time.Duration, proxyAddr ...string) (*http.Client, error) {
	transport := &http.Transport{}
	if len(proxyAddr) > 0 || proxyAddr[0] != "direct" {
		if strings.Contains(proxyAddr[0], "http") {
			proxyUrl, err := url.Parse(proxyAddr[0])
			if err != nil {
				return nil, err
			}
			// transport = &http.Transport{
			// 	Proxy: http.ProxyURL(proxyUrl),
			// }
			transport.Proxy = http.ProxyURL(proxyUrl)
		} else if strings.Contains(proxyAddr[0], "socks") {
			proxyUrl, err := url.Parse(proxyAddr[0])
			if err != nil {
				return nil, err
			}
			var proxyPass string
			if proxyPwd, ok := proxyUrl.User.Password(); ok {
				proxyPass = proxyPwd
			}
			dialer, err := proxy.SOCKS5("tcp", proxyUrl.Host, &proxy.Auth{
				User:     proxyUrl.User.Username(),
				Password: proxyPass,
			}, proxy.Direct)
			if err != nil {
				return nil, err
			}
			transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial("tcp", addr)
				},
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial("tcp", addr)
			}
		}
	}

	cli := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
	return cli, nil
}
