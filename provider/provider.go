package provider

import (
	"context"

	"github.com/hupe1980/notifier/provider/slack"
	"github.com/hupe1980/notifier/provider/sns"
	"github.com/hupe1980/notifier/provider/teams"
	"github.com/hupe1980/notifier/provider/webhook"
)

type Options struct {
	Slack   []*slack.Options   `mapstructure:"slack,omitempty"`
	SNS     []*sns.Options     `mapstructure:"sns,omitempty"`
	Teams   []*teams.Options   `mapstructure:"teams,omitempty"`
	Webhook []*webhook.Options `mapstructure:"webhook,omitempty"`
}

type Provider interface {
	Send(ctx context.Context, message string, extras map[string]string) error
}
