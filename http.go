package notifier

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type HTTPClientOptions struct {
	Proxy     string
	RateLimit int
}

func NewHTTPClient(options *HTTPClientOptions) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if options.Proxy != "" {
		proxyURL, err := url.Parse(options.Proxy)
		if err != nil {
			return nil, err
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	c := &http.Client{
		Transport: transport,
	}

	if options.RateLimit > 0 {
		c.Transport = newThrottledTransport(time.Second, options.RateLimit, transport)
	}

	return c, nil
}

type throttledTransport struct {
	roundTripperWrap http.RoundTripper
	ratelimiter      *rate.Limiter
}

func newThrottledTransport(limitPeriod time.Duration, requestCount int, transport http.RoundTripper) http.RoundTripper {
	return &throttledTransport{
		roundTripperWrap: transport,
		ratelimiter:      rate.NewLimiter(rate.Every(limitPeriod), requestCount),
	}
}

func (c *throttledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := c.ratelimiter.Wait(r.Context()); err != nil {
		return nil, err
	}

	return c.roundTripperWrap.RoundTrip(r)
}
