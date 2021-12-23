package sns

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/hupe1980/notifier/util"
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

	for _, opts := range pr.SNS {
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

	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithHTTPClient(pr.client),
		config.WithRegion(options.Region),
		config.WithSharedConfigProfile(options.Profile),
		config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
			aro.TokenProvider = stscreds.StdinTokenProvider
		}),
	)
	if err != nil {
		return err
	}

	snsClient := sns.NewFromConfig(awsCfg)

	input := &sns.PublishInput{
		Message:  &message,
		TopicArn: &options.TopicARN,
	}

	if _, err := snsClient.Publish(ctx, input); err != nil {
		return err
	}

	return nil
}
