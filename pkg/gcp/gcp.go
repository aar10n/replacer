package gcp

import (
	"errors"
	"io"
	"net/http"
	"time"
)

// InClusterProjectID returns the GCP project ID if the code is running inside a GKE cluster
func InClusterProjectID() (string, error) {
	client := &http.Client{}

	const url = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}

// WaitForWorkloadIdentity waits until the workload identity metadata server is ready.
// If the server is not ready within the specified timeout an error is returned.
func WaitForWorkloadIdentity(maxWait time.Duration) error {
	client := &http.Client{}

	const url = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	start := time.Now()
	for {
		resp, err := client.Do(req)
		if err != nil {
			goto retry
		}

		_, err = io.ReadAll(resp.Body)
		if err == nil {
			_ = resp.Body.Close()
			return nil
		}

	retry:
		if time.Since(start) > maxWait {
			_ = resp.Body.Close()
			return errors.New("timed out while waiting for metadata server")
		}

		time.Sleep(time.Second)
	}
}
