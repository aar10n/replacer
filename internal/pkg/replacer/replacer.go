package replacer

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aar10n/replacer/internal/pkg/providers"
	"github.com/aar10n/replacer/pkg/config"
)

const (
	replacerKeyPrefix = "replacer.agb.dev/"
	asyncThreshold    = 10 // prefetch async if more than 10 items
)

var (
	replacerRegexPattern = regexp.MustCompile(`<replace(?:\(([a-z_-]+)\))?:([^\n>]+)>`)
)

type replacement struct {
	full     string
	key      string
	provider *providers.Provider
}

// Replacer performs replacement on strings using various providers.
type Replacer struct {
	config    *Config
	rawConfig map[string]string
	providers map[string]*providers.Provider
}

// Config holds global replacer configuration options.
type Config struct {
	// Provider is the name of the default provider to use.
	Provider string `config:"provider"`
	// QuoteReplacements will quote replacement values.
	QuoteReplacements bool `config:"escape_replacements"`
	// IgnoreUnknownKeys will ignore unknown replacement keys.
	IgnoreUnknownKeys bool `config:"ignore_unknown_keys"`
}

// New creates a new replacer with the given config.
func New(cfg map[string]string) (*Replacer, error) {
	c := &Config{}
	err := config.LoadFromMapP(cfg, replacerKeyPrefix, c)
	if err != nil {
		return nil, err
	}

	r := &Replacer{
		config:    c,
		rawConfig: cfg,
		providers: make(map[string]*providers.Provider),
	}

	// instantiate default provider
	if c.Provider != "" {
		_, err := r.getProvider(c.Provider)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

// ReplaceAll replaces all replacement tags in the given string with
// values from the corresponding providers.
func (r *Replacer) ReplaceAll(s string) (string, error) {
	// collect all replacement candidates
	rkeys, err := r.getReplacementKeys(s)
	if err != nil {
		return "", err
	} else if rkeys == nil || len(rkeys) == 0 {
		return s, nil
	}

	// prefetch replacements
	err = r.prefetch(rkeys)
	if err != nil {
		return "", err
	}

	// replace all
	for _, rkey := range rkeys {
		val, err := rkey.provider.ValueFor(rkey.key)
		if err != nil {
			return "", err
		}

		s = strings.ReplaceAll(s, rkey.full, val)
	}

	return s, nil
}

func (r *Replacer) getProvider(name string) (*providers.Provider, error) {
	if p, ok := r.providers[name]; ok {
		return p, nil
	}

	// instantiate provider
	p, err := providers.Use(name)
	if err != nil {
		return nil, err
	}

	// load provider config
	err = config.LoadFromMapP(r.rawConfig, replacerKeyPrefix+name+".", p.ValueProvider)
	if err != nil {
		return nil, err
	}

	r.providers[name] = p
	return p, nil
}

func (r *Replacer) getReplacementKeys(s string) ([]replacement, error) {
	res := replacerRegexPattern.FindAllStringSubmatch(s, -1)
	if res == nil {
		return nil, nil
	}

	keys := make([]replacement, len(res))
	for i, match := range res {
		if match == nil || len(match) != 3 {
			continue
		}

		providerName := match[1]
		if providerName == "" {
			providerName = r.config.Provider
		}

		key := match[2]
		provider, err := r.getProvider(providerName)
		if err != nil {
			return nil, err
		}

		keys[i] = replacement{
			full:     match[0],
			key:      key,
			provider: provider,
		}
	}
	return keys, nil
}

func (r *Replacer) prefetch(rkeys []replacement) error {
	if len(rkeys) < asyncThreshold {
		// prefetch synchronously
		for _, rkey := range rkeys {
			_, err := rkey.provider.ValueFor(rkey.key)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// prefetch asynchronously
	wg := sync.WaitGroup{}
	errCh := make(chan error, len(rkeys))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	for _, rkey := range rkeys {
		wg.Add(1)
		go func(r replacement) {
			defer wg.Done()
			_, err := r.provider.ValueFor(r.key)
			if err != nil {
				cancel()
			}
			errCh <- err
		}(rkey)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
