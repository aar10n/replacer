package webhooks

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReplacerWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReplacerWebhook Suite")
}

var _ = Describe("ReplacerWebhook", func() {})
