package providers

import (
	"errors"
)

var (
	providers = make(map[string]Factory)
)

type Provider struct {
	Name string
	ValueProvider
}

type ValueProvider interface {
	// ValueFor returns a value for the given key. This is the core method of
	// the provider and is called whenever a string is being replaced.
	ValueFor(key string) (string, error)
}

type Closer interface {
	// Close should perform any cleanup required by the provider.
	Close()
}

type Factory func() (ValueProvider, error)

// Register associates a name with a provider factory. It should be called from
// the init function of the provider's package. This function panics if a provider
// with the same name is already registered.
func Register(name string, factory Factory) {
	if _, ok := providers[name]; ok {
		panic("provider already registered: " + name)
	}

	providers[name] = factory
}

// Use returns a new instance of the provider with the given name. It returns an
// error if no such provider is registered, or if the provider fails to initialize.
func Use(name string) (*Provider, error) {
	f, ok := providers[name]
	if !ok {
		if name == "" {
			return nil, errors.New("no provider given")
		}
		return nil, errors.New("provider not found: " + name)
	}

	vp, err := f()
	if err != nil {
		return nil, err
	}

	p := &Provider{
		Name:          name,
		ValueProvider: vp,
	}
	return p, nil
}

// Close performs cleanup required by the provider.
func (p *Provider) Close() {
	if closer, ok := p.ValueProvider.(Closer); ok {
		closer.Close()
	}
}
