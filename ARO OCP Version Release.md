## Making a Release
1. Make a new release branch based off the previous
2. Update /hack/update-aro-module-dependencies line 6 for the OCP version we're releasing.
3. Run /hack/update-go-module-dependencies.sh
    1. Make tea while we wait. When finished output will have something like:

``Update pkg/installer/generateconfig.go 's rhcosImage struct with:
SKU: "aro_420"
Version: "420.6.20251015",``

4. Update /pkg/installer/generateconfig.go with the output.
    1. SKU on line 168, 182
    2. Version on line 189

5. Verify the build and resolve dependency conflicts:
    1. Run `ARO_IMAGE=localhost/aro:test BUILDER_REGISTRY=registry.ci.openshift.org make image-aro` to test the build
    2. If you encounter dependency errors (undefined types, type mismatches, missing constants):
        - Check what version of dependencies the installer expects: `go mod graph | grep "installer.*<dependency-name>"`
        - Update the replace directives in go.mod to match the installer's expected versions
        - Common dependencies that may need updating: `github.com/openshift/api`, `github.com/openshift/client-go`
        - Re-run `go mod tidy && go mod vendor` after updating replace directives
        - Retry the build until it succeeds
    3. If errors persist, check the upstream installer's go.mod for the release branch to ensure version alignment. For example the K8s machinery should match the installer's at the same version.
     `` k8s.io/api v0.32.3                                                                                                                                                                            
        k8s.io/apimachinery v0.32.3                                                                                                                                                                   
        k8s.io/client-go v0.32.3 ``

## Test a cluster install
1. Go to https://oauth-openshift.apps.ci.l2s4.p1.openshiftapps.com/oauth/token/request
2. Click "RedHat_Internal_SSO" and log in
3. Copy the oc login --token=... command and run it in your terminal
4. Then run: oc registry login --to=$HOME/.docker/config.json 
5. oc registry login --to=$HOME/.docker/config.json
6.  ARO_IMAGE=localhost/aro:test BUILDER_REGISTRY=registry.ci.openshift.org make image-aro
