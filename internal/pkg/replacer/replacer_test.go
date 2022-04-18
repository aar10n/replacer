package replacer

import (
	"testing"

	"github.com/aar10n/replacer/internal/pkg/providers"
	_ "github.com/aar10n/replacer/internal/pkg/providers/gcp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReplacer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Replacer")
}

func makeTestProviderFactory(replacements map[string]string) providers.Factory {
	return func() (providers.ValueProvider, error) {
		p := providers.NewTestProvider(replacements)
		return p, nil
	}
}

var _ = Describe("Replacer", func() {
	Describe("ReplaceAll", func() {
		providers.Register("test",
			makeTestProviderFactory(map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			}),
		)

		It("should replace values with the default provider", func() {
			r, err := New(map[string]string{
				replacerKeyPrefix + "provider": "test",
			})
			Expect(err).ToNot(HaveOccurred())

			res, err := r.ReplaceAll(`
				hello <replace:key1> <replace:key2>
				<replace:key2> test
				more <replace:key3> <replace:key3>
				<replace:key4>
				yay
			`)
			expected := `
				hello value1 value2
				value2 test
				more value3 value3
				value4
				yay
			`

			Expect(res).To(Equal(expected))
		})

		It("should replace values with the specified provider", func() {
			r, err := New(map[string]string{})
			Expect(err).ToNot(HaveOccurred())

			res, err := r.ReplaceAll("<replace(test):key1>")
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("value1"))
		})

		It("should return an error if no provider is specified", func() {
			r, err := New(map[string]string{})
			Expect(err).ToNot(HaveOccurred())

			_, err = r.ReplaceAll("<replace:key1>")
			Expect(err).To(HaveOccurred())
		})
	})

	//	Describe("ReplaceAll (gcp)", func() {
	//		It("should replace values with the gcp provider", func() {
	//			r, err := New(map[string]string{
	//				replacerKeyPrefix + "provider":       "gcp",
	//				replacerKeyPrefix + "gcp.project_id": "cohere-cd",
	//			})
	//			Expect(err).ToNot(HaveOccurred())
	//
	//			res, err := r.ReplaceAll(`
	//aurthur_secrets.yaml
	//	<replace:production-aurthur-secrets>
	//admin_password: <replace:projects/cohere-cd/secrets/production-admin-password>
	//aurthur_admin_password: <replace:production-aurthur-admin-password>
	//			`)
	//			Expect(err).ToNot(HaveOccurred())
	//			fmt.Println("res:", res)
	//		})
	//	})
})
