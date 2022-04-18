package gcp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aar10n/replacer/internal/pkg/providers"
	"github.com/aar10n/replacer/pkg/cache"
	"github.com/aar10n/replacer/pkg/gcp"
)

// GCP Secret Manager Provider

var (
	fullPathPattern  = regexp.MustCompile(`^projects/([\w-]+)/secrets/([\w-]+)(?:/versions/([\d]+|latest))?$`)
	shortPathPattern = regexp.MustCompile(`^(?:([\w-]+)/)?([\w-]+)$`)
)

type SecretManagerProvider struct {
	client *gcp.SecretManagerClient
	cache  *cache.Cache

	ProjectID string `config:"project_id"`
}

func SecretManagerProviderFactory() (providers.ValueProvider, error) {
	client, err := gcp.NewSecretManagerClient()
	if err != nil {
		return nil, err
	}

	p := &SecretManagerProvider{
		client: client,
		cache:  cache.New(),
	}

	return p, nil
}

func (p *SecretManagerProvider) ValueFor(key string) (string, error) {
	newKey, err := p.getSecretPath(key)
	if err != nil {
		return "", err
	}

	key = newKey
	if value := p.cache.Get(key); value != nil {
		return value.(string), nil
	}

	value, err := p.client.GetSecret(key)
	if err != nil {
		return "", err
	}

	p.cache.Set(key, value)
	return value, nil
}

func (p *SecretManagerProvider) Close() {
	p.client.Close()
}

func (p *SecretManagerProvider) getSecretPath(key string) (string, error) {
	key = strings.Trim(key, " \t")

	if res := fullPathPattern.FindStringSubmatch(key); res != nil {
		if res[3] == "" {
			return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", res[1], res[2]), nil
		}
		return key, nil
	} else if res = shortPathPattern.FindStringSubmatch(key); res != nil {
		if res[1] == "" {
			if p.ProjectID == "" {
				return "", fmt.Errorf("missing project_id in path or config")
			}
			return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", p.ProjectID, res[2]), nil
		}
		return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", res[1], res[2]), nil
	}
	return "", fmt.Errorf("invalid secret path: %s", key)
}

//

func init() {
	// register the provider
	providers.Register("gcp", SecretManagerProviderFactory)
}
