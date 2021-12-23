package sns

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"go.uber.org/multierr"
)

type Options struct {
	ID       string `mapstructure:"id,omitempty"`
	Region   string `mapstructure:"region,omitempty"`
	Profile  string `mapstructure:"profile,omitempty"`
	TopicARN string `mapstructure:"topicArn,omitempty"`
	Template string `mapstructure:"template,omitempty"`
}

type Provider struct {
	SNS    []*Options `mapstructure:"sns,omitempty"`
	client *http.Client
}

func New(client *http.Client, options []*Options) (*Provider, error) {
	provider := &Provider{
		client: client,
	}

	provider.SNS = append(provider.SNS, options...)

	return provider, nil
}

func (pr *Provider) Send(ctx context.Context, message string, extras map[string]string) error {
	var sendErr error

	for _, p := range pr.SNS {
		awsCfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithHTTPClient(pr.client),
			config.WithRegion(p.Region),
			config.WithSharedConfigProfile(p.Profile),
			config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
				aro.TokenProvider = stscreds.StdinTokenProvider
			}),
		)
		if err != nil {
			sendErr = multierr.Append(sendErr, err)
			continue
		}

		snsClient := sns.NewFromConfig(awsCfg)

		input := &sns.PublishInput{
			Message:  &message,
			TopicArn: &p.TopicARN,
		}

		if _, err := snsClient.Publish(ctx, input); err != nil {
			sendErr = multierr.Append(sendErr, err)
			continue
		}
	}

	return sendErr
}
