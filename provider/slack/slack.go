package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hupe1980/notifier/util"
	"go.uber.org/multierr"
)

const Name = "slack"

type Options struct {
	ID         string `yaml:"id,omitempty"`
	WebhookURL string `yaml:"webhookUrl,omitempty"`
	Template   string `yaml:"template,omitempty"`
}

type Provider struct {
	Slack  []*Options `yaml:"slack,omitempty"`
	client *http.Client
}

func New(client *http.Client, options []*Options) (*Provider, error) {
	provider := &Provider{
		client: client,
	}

	provider.Slack = append(provider.Slack, options...)

	return provider, nil
}

func (pr *Provider) Send(ctx context.Context, message string, extras map[string]string) error {
	var sendErr error

	for _, opts := range pr.Slack {
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

	payload, err := json.Marshal(WebhookRequest{
		Text: message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", options.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pr.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return err
		}

		return fmt.Errorf("error while sending slack message: %s(%d)", buf.String(), resp.StatusCode)
	}

	return nil
}
