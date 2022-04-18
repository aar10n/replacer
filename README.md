# Kubernetes Replacer

A Kubernetes webhook that performs replacements on Secrets and ConfigMaps using one of the 
supported "providers". It simplifies secret management by allowing you to inject secrets into 
your application without the use of custom resources. It also allows you to compose secrets 
from multiple sources, and in whatever format you want.

## Why?

This project was heavily inspired by the tool [ArgoCD Vault Plugin](https://github.com/argoproj-labs/argocd-vault-plugin).
Unlike many other secret management tools, it performs replacements the yaml before it is applied 
to the cluster. This has the benefit of requiring no custom resources or controllers, and it allows 
you to compose and combine multiple secrets into a single resource. The downside is that it requires
you to install the application into your CI/CD pipeline, and it makes testing locally less convenient.

This project was created as a way to get the functionality of ArgoCD Vault Plugin into Kubernetes. 
Through the use of a mutating webhook, it performs the similar replacement functionality as ArgoCD 
Vault Plugin, but it does so as part of the normal `kubectl apply` process.

## Example

The following example uses the [GoogleSecretManager](#GoogleSecretManager) provider:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  annotations:
    replacer.agb.dev/provider: gcp
stringData:
  key1: <replace:my-project/some-secret>
  secrets.yaml: |
    <replace:my-project/some-yaml-secret>
    api_token: <replace:my-project/some-token-secret>
```

## Providers

A provider is a backend that provides replacements for keys inside of `<replace:>` templates. 
You can select a default provider with the `replacer.agb.dev/provider` annotation on your resource,
or with the `<replace(<provider>):>` template syntax. 

Currently, only the `gcp` provider is supported, but it is very easy to add a new provider and
pull requests are welcome.

#### Provider Configuration

All provider-specific configuration options are specified via annotations with the
prefix `replacer.agb.dev/` followed by the provider name, a period, and finally the 
key. If the type of a parameter is listed as `boolean` or `integer`, the value should 
be a string representation of the value (e.g. `replacer.agb.dev/key: "true"`).

### GoogleSecretManager

Provider for Google Cloud Platform's [Secret Manager](https://cloud.google.com/secret-manager).

#### Usage

```yaml
metadata:
  annotations:
    replacer.agb.dev/provider: gcp
```

The `gcp` provider accepts both the full resource path to the secret, or shorter forms which 
include just the secret name and project (not nessecary if a default project is given). The
default version used in all cases where is it not specified is `latest`.

Secret path examples:
  * `<replace:projects/my-project/secrets/my-secret>` 
  * `<replace:projects/my-project/secrets/my-secret/versions/latest>`
  * `<replace:projects/my-project/secrets/my-secret/versions/1>`
  * `<replace:my-project/my-secret>`
  * `<replace:my-secret>` (only if the `project_id` option is provided)

#### Configuration

| Key          | Type   | Description                                       |
|--------------|--------|---------------------------------------------------|
| `project_id` | string | The default project id to use when none is given. |

### AWSSecretManager

Planned


## License

MIT License, see the LICENSE file.

