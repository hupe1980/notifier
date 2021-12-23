package main

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/hupe1980/notifier"
	"github.com/hupe1980/notifier/provider"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

const (
	version = "dev"
)

func main() {
	var opts struct {
		config    string
		proxy     string
		rateLimit int
		bulk      bool
	}

	rootCmd := &cobra.Command{
		Use:     "notifier [data]",
		Version: version,
		Short:   "Tiny helper for publishing notifications on a variety of supported platforms",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := readConfig(opts.config)
			if err != nil {
				return err
			}

			httpclientOptions := &notifier.HTTPClientOptions{
				Proxy: opts.proxy,
			}

			if opts.proxy != "" {
				httpclientOptions.Proxy = opts.proxy
			}

			if opts.rateLimit != 0 {
				httpclientOptions.RateLimit = opts.rateLimit
			}

			c, err := notifier.NewHTTPClient(httpclientOptions)
			if err != nil {
				return err
			}

			n, err := notifier.New(c, cfg.Providers)
			if err != nil {
				return err
			}

			filename := ""
			if len(args) == 1 {
				filename = args[0]
			}

			in, err := notifier.NewInput(filename)
			if err != nil {
				return err
			}
			defer in.Close()

			extras, err := additionalInfos()
			if err != nil {
				return err
			}

			if opts.bulk {
				bulk, err := in.Bulk()
				if err != nil {
					return nil
				}

				if err := n.Send(context.Background(), bulk, extras); err != nil {
					return err
				}

				return nil
			}

			for line := range in.Line() {
				if err := n.Send(context.Background(), line, extras); err != nil {
					return err
				}
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&opts.config, "config", "c", "", "path to notifier configuration file (default: $HOME/.config/notifier/config.yaml)")
	rootCmd.Flags().StringVarP(&opts.proxy, "proxy", "", "", "proxy url")
	rootCmd.Flags().IntVarP(&opts.rateLimit, "rate-limit", "", 0, "maximum number of HTTP requests per second")
	rootCmd.Flags().BoolVarP(&opts.bulk, "bulk", "b", false, "enable bulk processing")

	if err := rootCmd.Execute(); err != nil {
		for _, v := range multierr.Errors(err) {
			fmt.Fprintln(os.Stderr, v)
		}
		os.Exit(1)
	}
}

type config struct {
	Providers *provider.ProviderOptions `mapstructure:"providers,omitempty"`
	Proxy     string                    `mapstructure:"proxy,omitempty"`
	RateLimit int                       `mapstructure:"rateLimit,omitempty"`
}

func readConfig(filename string) (*config, error) {
	if filename != "" {
		viper.SetConfigFile(filename)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(fmt.Sprintf("%s/.config/notifier", home))
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func additionalInfos() (map[string]string, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Username": user.Username,
		"Hostname": hostname,
	}, nil
}
