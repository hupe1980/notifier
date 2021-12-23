# notifier
> Tiny helper for publishing notifications on a variety of supported platforms: 
- #slack
- AWS SNS
- Teams
- Custom Webhooks

```bash
nmap -p80,443 scanme.nmap.org | notifier -b
```

<div align="left">
  <img src="slack.png" alt="slack" width="700px"></a>
</div>

## Installing
You can install the pre-compiled binary in several different ways

### homebrew tap:
```bash
brew tap hupe1980/notifier
brew install notifier
```

### snapcraft:
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/notifier)
```bash
sudo snap install notifier
```

### scoop:
```bash
scoop bucket add notifier https://github.com/hupe1980/notifier-bucket.git
scoop install notifier
```

### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/notifier/releases) and install them with the appropriate tools.

### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/notifier/releases) and copy to the desired location.

## Usage
```console
Usage:
  notifier [data] [flags]

Flags:
  -b, --bulk             enable bulk processing
  -c, --config string    path to notifier configuration file (default: $HOME/.config/notifier/config.yaml)
  -h, --help             help for notifier
      --proxy string     proxy url
      --rate-limit int   maximum number of HTTP requests per second
  -v, --version          version for notifier
```

### Provider
```yaml
proxy: http://proxy.org
rateLimit: 5
providers:
  webhook:
    - id: webhook
      url: https://webhook.org
      method: POST
      template: '{{ .Message }}'
      headers:
        Content-Type: application/json
        X-Api-Key: 4711
  slack:
    - id: slack
      webhookUrl: https://hooks.slack.com/services/xxx
      template: '{{ .Message }}'
```

### Template
The following placeholders are supported:
- {{ .Message }}
- {{ .Username }}
- {{ .Hostname }}
## License
[MIT](LICENCE)

