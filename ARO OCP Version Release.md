## Making a Release
1. Make a new release branch based off the previous
2. Update /hack/update-aro-module-dependencies line 6 for the OCP version we’re releasing.
3. Run /hack/update-go-module-dependencies.sh
    1. Make tea while we wait. When finished output will have something like:

``Update pkg/installer/generateconfig.go 's rhcosImage struct with:
SKU: "aro_420"
Version: "9.6.20251015",``

4. Update /pkg/installer/generateconfig.go with the output.
    1. SKU on line 168, 182
    2. Version on line 189
