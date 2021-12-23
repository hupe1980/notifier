package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hupe1980/notifier/util"
	"go.uber.org/multierr"
)

type Options struct {
	ID         string `mapstructure:"id,omitempty"`
	WebhookURL string `mapstructure:"webhookUrl,omitempty"`
	Title      string `mapstructure:"title,omitempty"`
	ThemeColor string `mapstructure:"themeColor,omitempty"`
	Template   string `mapstructure:"template,omitempty"`
}

type Provider struct {
	Teams  []*Options `mapstructure:"teams,omitempty"`
	client *http.Client
}

func New(client *http.Client, options []*Options) (*Provider, error) {
	provider := &Provider{
		client: client,
	}

	provider.Teams = append(provider.Teams, options...)

	return provider, nil
}

func (pr *Provider) Send(ctx context.Context, message string, extras map[string]string) error {
	var sendErr error

	for _, opts := range pr.Teams {
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

	sections := []Section{}
	for _, line := range strings.Split(message, "\n") {
		sections = append(sections, Section{
			Text: line,
		})
	}

	summary := options.Title
	if summary == "" && len(sections) > 0 {
		summary = sections[0].Text
		if len(summary) > 20 {
			summary = summary[:21]
		}
	}

	payload, err := json.Marshal(MessageCard{
		CardType:   "MessageCard",
		Context:    "http://schema.org/extensions",
		Markdown:   true,
		Title:      options.Title,
		ThemeColor: options.ThemeColor,
		Summary:    summary,
		Sections:   sections,
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

	resp.Body.Close()

	return nil
}
