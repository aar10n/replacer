package webhooks

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func RegisterWebhooksWithManager(mgr ctrl.Manager) error {
	server := mgr.GetWebhookServer()
	server.Register("/replace", &webhook.Admission{
		Handler: &ReplacerWebhook{Client: mgr.GetClient()},
	})
	return nil
}
