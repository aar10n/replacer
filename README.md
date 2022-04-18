# Secret Replacer
A kubernetes mutating webhook to replace values in resources with secrets fetched from a backend.


### Installing on a private cluster
When installing on a private cluster you need to create a firewall rule to allow the master node to access the webhook.

See: https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#add_firewall_rules
