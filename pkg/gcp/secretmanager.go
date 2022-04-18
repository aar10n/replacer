package gcp

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManagerClient struct {
	client *secretmanager.Client
}

// NewSecretManagerClient opens a connection to the secret manager service.
func NewSecretManagerClient() (*SecretManagerClient, error) {
	client, err := secretmanager.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	return &SecretManagerClient{client: client}, nil
}

// GetSecret accesses the given secret and returns its payload.
// The secret must be in the following format:
//   projects/<project-id>/secrets/<secret-name>/versions/<version-id>
func (s *SecretManagerClient) GetSecret(secret string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secret,
	}
	resp, err := s.client.AccessSecretVersion(context.TODO(), req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload.Data), nil
}

// Close closes the connection to the secret manager service.
func (s *SecretManagerClient) Close() {
	_ = s.client.Close()
}
