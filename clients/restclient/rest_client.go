package restclient

import (
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks = make(map[string]*Mock)
)

type Mock struct {
	Url string
	HttpMethod string
	Response *http.Response
	Err error
}

func getMockId(httpMethod string, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func StartMockups() {
	enabledMocks = true
}

func StopMockups() {
	enabledMocks = false
}

func FlushMocks() {
	mocks = make(map[string]*Mock)
}

func AddMockup(mock Mock) {
	mocks[getMockId(mock.HttpMethod, mock.Url)] = &mock
}

type HttpGetter interface {
	Get(string, http.Header) (*http.Response, error)
}

type Config struct {
	BaseUrl string
}

type client struct {
	baseUrl string
	c *http.Client
}

var _ HttpGetter = (*client)(nil)

func New(cfg Config) *client {
	return &client{
		baseUrl: cfg.BaseUrl,
		c: &http.Client{},
	}
}

func (c *client) Client() *http.Client  {
	return c.c
}

func (c *client) Get(uri string, header http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseUrl, uri)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		request.Header = header
	}

	return c.c.Do(request)
}
