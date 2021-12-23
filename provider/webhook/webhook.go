package webhook

import (
	"bytes"
	"context"
	"net/http"

	"github.com/hupe1980/notifier/util"
	"go.uber.org/multierr"
)

const Name = "webhook"

type Options struct {
	ID       string            `mapstructure:"id,omitempty"`
	URL      string            `mapstructure:"url,omitempty"`
	Method   string            `mapstructure:"method,omitempty"`
	Headers  map[string]string `mapstructure:"headers,omitempty"`
	Template string            `mapstructure:"template,omitempty"`
}

type Provider struct {
	Webhook []*Options `mapstructure:"webhook,omitempty"`
	client  *http.Client
}

func New(client *http.Client, options []*Options) (*Provider, error) {
	provider := &Provider{
		client: client,
	}

	provider.Webhook = append(provider.Webhook, options...)

	return provider, nil
}

func (pr *Provider) Send(ctx context.Context, message string, extras map[string]string) error {
	var sendErr error

	for _, opts := range pr.Webhook {
		err := pr.send(ctx, message, extras, opts)
		if err != nil {
			sendErr = multierr.Append(sendErr, err)
		}
	}

	return sendErr
}

func (pr *Provider) send(ctx context.Context, message string, extras map[string]string, options *Options) error {
	message, err := util.ExecuteTemplate(options.ID, options.Template, message, extras)
	if err != nil {
		return err
	}

	body := bytes.NewBufferString(message)

	req, err := http.NewRequestWithContext(ctx, options.Method, options.URL, body)
	if err != nil {
		return err
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	resp, err := pr.client.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}
