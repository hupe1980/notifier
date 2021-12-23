package notifier

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hupe1980/notifier/provider"
	"github.com/hupe1980/notifier/provider/slack"
	"github.com/hupe1980/notifier/provider/sns"
	"github.com/hupe1980/notifier/provider/teams"
	"github.com/hupe1980/notifier/provider/webhook"
	"github.com/hupe1980/notifier/util"
	"go.uber.org/multierr"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var ansiRE = regexp.MustCompile(ansi)

type Notifier struct {
	providers []provider.Provider
	client    *http.Client
}

func New(client *http.Client, providers []string, options *provider.Options) (*Notifier, error) {
	notifier := &Notifier{
		client: client,
	}

	if options.Slack != nil && util.ContainsOrIsEmpty(providers, slack.Name) {
		provider, err := slack.New(client, options.Slack)
		if err != nil {
			return nil, err
		}

		notifier.providers = append(notifier.providers, provider)
	}

	if options.SNS != nil && util.ContainsOrIsEmpty(providers, sns.Name) {
		provider, err := sns.New(client, options.SNS)
		if err != nil {
			return nil, err
		}

		notifier.providers = append(notifier.providers, provider)
	}

	if options.Teams != nil && util.ContainsOrIsEmpty(providers, teams.Name) {
		provider, err := teams.New(client, options.Teams)
		if err != nil {
			return nil, err
		}

		notifier.providers = append(notifier.providers, provider)
	}

	if options.Webhook != nil && util.ContainsOrIsEmpty(providers, webhook.Name) {
		provider, err := webhook.New(client, options.Webhook)
		if err != nil {
			return nil, err
		}

		notifier.providers = append(notifier.providers, provider)
	}

	return notifier, nil
}

func (n *Notifier) Send(ctx context.Context, message string, extras map[string]string) error {
	var sendErr error

	message = ansiRE.ReplaceAllString(message, "")

	if message == "" {
		return nil
	}

	for _, p := range n.providers {
		if err := p.Send(ctx, message, extras); err != nil {
			sendErr = multierr.Append(sendErr, err)
		}
	}

	return sendErr
}
