# Test paragraph

You can build and publish to your personal Quay account installer-aro-wrapper image with as follow:
```
export ARO_IMAGE=quay.io/<your_account>/aroinstaller:<image_tag>
make publish-image-aro
```

For the build step to succeed, we are now pulling images from quay.io/openshift-release-dev and they aren't publicly available. So we need to have extra auth for docker/pod
man. Needed information can be found in `secrets/env` in ARO-RP.
Note `secrets` repository isn't commited and is retrieved with `make secrets` command (details in ARO-RP documentation).
The needed token is the one used for 'quay.io' in the `PULL_SECRET` environment variable.

To be able to pull images for build and push your image afterward to your own account, `auth.json` should be something like
```
{
  "auths": {
    "arointsvc.azurecr.io": {
      "auth": "<redacted_token>"
    },
      "quay.io/<your_account>": {
        "auth": "<redacted_token>"
    },
      "quay.io/openshift-release-dev": {
        "auth": "<redacted_token_from_ARO-rp>"
    }
  }
}
```
