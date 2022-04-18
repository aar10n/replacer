package webhooks

import (
	"context"
	"net/http"

	_ "github.com/aar10n/replacer/internal/pkg/providers/gcp"
	"github.com/aar10n/replacer/internal/pkg/replacer"

	"gomodules.xyz/jsonpatch/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type ReplacerWebhook struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (w *ReplacerWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := logf.FromContext(ctx)

	var patches []jsonpatch.Operation
	var err error
	switch req.RequestKind.Kind {
	case "Secret":
		secret := &corev1.Secret{}
		err = w.decoder.Decode(req, secret)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		} else if len(secret.Data) == 0 {
			return admission.Allowed("no data to replace")
		}

		log.Info("Handle.Secret", "name", secret.Name, "namespace", secret.Namespace)
		patches, err = w.replaceInSecret(ctx, secret)
	case "ConfigMap":
		cm := &corev1.ConfigMap{}
		err = w.decoder.Decode(req, cm)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		} else if len(cm.Data) == 0 {
			return admission.Allowed("no data to replace")
		}

		log.Info("Handle.ConfigMap", "name", cm.Name, "namespace", cm.Namespace)
		patches, err = w.replaceInConfigMap(ctx, cm)
	default:
		return admission.Allowed("not a secret or configmap")
	}

	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	} else if patches == nil || len(patches) == 0 {
		return admission.Allowed("no changes")
	}
	return admission.Patched("replacer changes", patches...)
}

func (w *ReplacerWebhook) InjectDecoder(d *admission.Decoder) error {
	w.decoder = d
	return nil
}

//

func (w *ReplacerWebhook) replaceInSecret(ctx context.Context, secret *corev1.Secret) ([]jsonpatch.Operation, error) {
	r, err := replacer.New(secret.Annotations)
	if err != nil {
		return nil, err
	}

	var patches []jsonpatch.Operation
	for k, oldValueB := range secret.Data {
		oldValue := string(oldValueB)
		newValue, err := r.ReplaceAll(oldValue)
		if err != nil {
			return nil, err
		} else if newValue == oldValue {
			continue
		}

		patches = append(patches, jsonpatch.Operation{
			Operation: "replace",
			Path:      "/data/" + k,
			Value:     newValue,
		})
	}

	return patches, nil
}

func (w *ReplacerWebhook) replaceInConfigMap(ctx context.Context, cm *corev1.ConfigMap) ([]jsonpatch.Operation, error) {
	r, err := replacer.New(cm.Annotations)
	if err != nil {
		return nil, err
	}

	var patches []jsonpatch.Operation
	for k, v := range cm.Data {
		newV, err := r.ReplaceAll(v)
		if err != nil {
			return nil, err
		}

		if newV != v {
			patches = append(patches, jsonpatch.Operation{
				Operation: "replace",
				Path:      "/data/" + k,
				Value:     newV,
			})
		}
	}

	return []jsonpatch.Operation{}, nil
}
