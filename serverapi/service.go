package serverapi

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Service app store server api 实现
// api: https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
type Service struct {
	client *http.Client
	// debug 将打印请求日志，用于开发测试
	debug bool

	// 是否 sandbox 环境
	sandbox bool
	// api token
	token *Token
}

func NewService(token *Token) *Service {
	return &Service{
		client: http.DefaultClient,
		token:  token,
	}
}

func (s *Service) clone() *Service {
	ns := new(Service)
	*ns = *s

	return ns
}

func (s *Service) Debug(debug bool) *Service {
	ns := s.clone()
	ns.debug = debug
	return ns
}

func (s *Service) Sandbox(sandbox bool) *Service {
	ns := s.clone()
	ns.sandbox = sandbox
	return ns
}

func (s *Service) Host() string {
	if s.sandbox {
		return "https://api.storekit-sandbox.itunes.apple.com"
	}
	return "https://api.storekit.itunes.apple.com"
}

func (s *Service) get(ctx context.Context, path string, query url.Values) (int, []byte, error) {
	u, err := url.Parse(path)
	if err != nil {
		return 0, nil, err
	}

	u.RawQuery = query.Encode()
	return s.request(ctx, "GET", u.String(), nil)
}

func (s *Service) request(ctx context.Context, method string, path string, body io.Reader) (int, []byte, error) {
	token, err := s.token.Get()
	if err != nil {
		return 0, nil, err
	}

	timeout := s.token.conf.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method,
		fmt.Sprintf("%s/%s", s.Host(), strings.TrimPrefix(path, "/")), body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	if s.debug {
		reqLog, _ := httputil.DumpRequestOut(req, true)
		log.Printf("[debug] %s", reqLog)
	}

	start := time.Now()
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if s.debug {
		respLog, _ := httputil.DumpResponse(resp, true)
		log.Printf("[debug] [latency:%s] %s", time.Now().Sub(start), respLog)
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	// 响应错误
	err = handleApiErr(resp.StatusCode, payload)
	return resp.StatusCode, payload, err
}
