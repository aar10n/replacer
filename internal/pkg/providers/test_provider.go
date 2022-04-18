package providers

import "errors"

type TestProvider struct {
	replacements map[string]string
}

func NewTestProvider(replacements map[string]string) *TestProvider {
	return &TestProvider{
		replacements: replacements,
	}
}

func (p *TestProvider) ValueFor(key string) (string, error) {
	if value, ok := p.replacements[key]; ok {
		return value, nil
	}
	return "", errors.New("key not found")
}
