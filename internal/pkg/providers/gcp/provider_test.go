package gcp

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GCP Provider Suite")
}

var _ = Describe("GCP Provider", func() {
	Describe("getSecretPath", func() {
		It("should return a correct path for a short-form secret with project id", func() {
			provider := &SecretManagerProvider{ProjectID: ""}
			path := "my-project/secret-name"

			result, err := provider.getSecretPath(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("projects/my-project/secrets/secret-name/versions/latest"))
		})
		It("should return a correct path for a short-form secret with default project id", func() {
			provider := &SecretManagerProvider{ProjectID: "my-project"}
			path := "secret-name"

			result, err := provider.getSecretPath(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("projects/my-project/secrets/secret-name/versions/latest"))
		})
		It("should return a correct path for a long-form secret without version", func() {
			provider := &SecretManagerProvider{}
			path := "projects/my-project/secrets/secret-name"

			result, err := provider.getSecretPath(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("projects/my-project/secrets/secret-name/versions/latest"))
		})
		It("should return a correct path for a long-form secret with version", func() {
			provider := &SecretManagerProvider{}
			path := "projects/my-project/secrets/secret-name/versions/1"

			result, err := provider.getSecretPath(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("projects/my-project/secrets/secret-name/versions/1"))
		})
		It("should return an error for a short-form path without an explicit or default project id", func() {
			provider := &SecretManagerProvider{}
			path := "secret-name"

			result, err := provider.getSecretPath(path)
			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(""))
		})
	})
})
