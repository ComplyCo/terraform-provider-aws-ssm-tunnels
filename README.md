# AWS SSM Tunnels provider

This provider's purpose is to make it possible to use AWS SSM port forwarding to connect to private resources (such as EKS, RDS, ...) from Terraform Cloud. This provider was inspired by https://github.com/flaupretre/terraform-ssh-tunnel. That module works on machines which can have the AWS SSM Session Manager Plugin installed. This is not the case in Terraform Cloud. As such, we need a way to bundle the Session Manager Plugin as a provider. Luckily for us, it is available as Go code (https://github.com/aws/session-manager-plugin). Please note, that this provider does not run the latest version because of this open issue (https://github.com/aws/session-manager-plugin/issues/73) where the UUID library in Session Manager Plugin has breaking changes.

## Challenges and limitations

Building this provider was a good learning experience for Terraform's framework of managing resources. The framework doesn't seem to have a concept of a provider "wrapper" were one provider can exist for the duration of another provider (in this case, the tunnel needs to stay open while the kubernetes, helm, postgres, ... provider is doing things with resources).

In developing this provider, we saw that Terraform shut down the provider when the provider was "done" with its last resource. If we only configured `data.awsssmtunnels_remote_tunnel.rds` and no other resource, the provider would shut down as soon as the data resource returned with the local host+port for the tunnel. This would then shut down the tunnel before the kubernetes/helm/postgres/.... providers could use it.

To get around this we added the `data.awsssmtunnels_keepalive.rds` resource which requires the caller to pass in all resources for provider using the tunnel to a `depends_on` lifecycle hook. This is a pretty poor developer experience, but it was all we could come up with at the present for keeping the tunnel running until all the resources that needed it were finished using the tunnel.

## Quality of the code

This provider is in an early-development state and has room for API, documentation, and testing improvements.

Please note, this was scaffolded from the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and there may be code/examples/configuration that needs cleanup.
