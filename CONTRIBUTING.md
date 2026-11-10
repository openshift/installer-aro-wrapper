# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution. For
details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only need
to do this once across all repositories using our CLA.

This project has adopted the [Microsoft Open Source Code of
Conduct](https://opensource.microsoft.com/codeofconduct/). For more information
see the [Code of Conduct
FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact
[opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional
questions or comments.

## Testing AKS Workload Identity

Production installer Jobs can use a dedicated hosting identity instead of the
legacy Hive AKS agent-pool identity. This identity is independent of the customer
cluster's service principal or managed identities; both target identity models
continue to use their existing installer inputs.

The RP supplies the following non-secret deployment configuration:

| Variable | Value |
|---|---|
| `ARO_INSTALLER_IDENTITY_CLIENT_ID` | Expected installer hosting UAMI client ID |
| `AZURE_ENVIRONMENT` | `AzurePublicCloud` or `AzureUSGovernmentCloud` |
| `ARO_AZURE_SUBSCRIPTION_ID` | RP deployment subscription, not the customer subscription |
| `ARO_LOCATION` | RP deployment region |
| `ARO_RESOURCEGROUP` | RP deployment resource group |
| `ARO_KEYVAULT_PREFIX` | RP Key Vault prefix, used to locate the `-svc` vault |

Existing production settings, including `ARO_AZURE_FP_CLIENT_ID`,
`ARO_DOMAIN_NAME`, `ARO_GATEWAY_DOMAINS`, `ARO_GATEWAY_RESOURCEGROUP`, and
`ARO_ACR_RESOURCE_ID`, are still required.

Use namespace `aro-installer` and ServiceAccount `aro-installer`, annotated with
`azure.workload.identity/client-id` for the hosting identity. Label the pod
`azure.workload.identity/use: "true"`. Its federated credential must match the
cluster's OIDC issuer, subject
`system:serviceaccount:aro-installer:aro-installer`, and audience
`api://AzureADTokenExchange`. The AKS webhook supplies `AZURE_CLIENT_ID`,
`AZURE_TENANT_ID`, and `AZURE_FEDERATED_TOKEN_FILE`, including the projected token
volume. Do not put tokens or client secrets in ordinary environment variables.

The wrapper uses `WorkloadIdentityCredential` for hosting-service access and
explicit deployment metadata without contacting VM IMDS. The expected hosting
client ID, if supplied, must match the webhook's client ID. Missing or empty
configuration fails instead of falling back to the node identity. Without any
hosting-identity or webhook variables, the existing Hive path is unchanged.
`ARO_RP_MODE=development` retains the existing local credential and proxy path.

Run the focused regression tests with:

```sh
go test -tags=altinfra,aro,containers_image_openpgp ./pkg/util/instancemetadata ./pkg/util/azureclient ./pkg/env -count=1
```

To test an unmerged wrapper branch, build and push a uniquely tagged image to an
authorized non-production registry:

```sh
make image-aro ARO_IMAGE=<registry>/<repository>:<unique-branch-tag>
docker push <registry>/<repository>:<unique-branch-tag>
```

Configure the test RP's selected `OpenShiftVersion` entry's `installerPullspec`
to that image (preferably its digest), preserving the matching OpenShift release
pullspec. Use the local ARO-RP branch with `aksjob` selected and explicit cloud,
deployment metadata, and installer hosting client ID wired into its Job inputs.
Do not replace shared release image tags.

Before a real install, verify image pulling, token projection, Entra exchange,
access to a designated non-production Key Vault value, ARM access, networking,
and scoped RP Kubernetes permissions with a synthetic workload. Then validate
the wrapper input mounts without creating infrastructure. Finally test both
customer identity models through create/validate/delete, RP restart reattachment,
failure diagnostics, and execution-scoped cleanup. Unit tests do not establish
that these infrastructure prerequisites are deployed or authorized.
