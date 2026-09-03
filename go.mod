module github.com/openshift/installer-aro-wrapper

go 1.25.8

require (
	github.com/Azure/ARO-RP/pkg/api v0.0.0-20260625191944-7d3e42bd8ad5 // v20260625.03
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7 v7.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.8.1
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.2
	github.com/Azure/go-autorest/autorest v0.11.30
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/codahale/etm v0.0.0-20141003032925-c00c9e6fb4c9
	github.com/coreos/go-semver v0.3.1
	github.com/coreos/ignition/v2 v2.24.0
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang/mock v1.7.0-rc.1
	github.com/google/go-cmp v0.7.0
	github.com/jongio/azidext/go/azidext v0.5.0
	github.com/onsi/gomega v1.38.2
	github.com/openshift/api v0.0.0-20260817165856-9fda5179a7f7
	github.com/openshift/installer v1.4.22-ec5.0.20260820153115-99421e0f93f0
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af
	github.com/spf13/cobra v1.10.1
	github.com/stretchr/testify v1.11.1
	github.com/vincent-petithory/dataurl v1.0.0
	go.uber.org/mock v0.6.0
	golang.org/x/crypto v0.52.0
	k8s.io/api v0.35.1
	k8s.io/apimachinery v0.35.1
	k8s.io/client-go v0.35.1
	sigs.k8s.io/cluster-api-provider-azure v1.21.1-0.20250929163617-2c4eaa611a39
	sigs.k8s.io/yaml v1.6.0
)

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.121.6 // indirect
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/kms v1.22.0 // indirect
	cloud.google.com/go/longrunning v0.6.7 // indirect
	cloud.google.com/go/monitoring v1.24.2 // indirect
	cloud.google.com/go/storage v1.57.0 // indirect
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3 v3.0.0-beta.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.5.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.24 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.1 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.2 // indirect
	github.com/Azure/go-autorest/tracing v0.6.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.6.0 // indirect
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.30.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 // indirect
	github.com/IBM-Cloud/bluemix-go v0.0.0-20231026045515-5c31bb9104f7 // indirect
	github.com/IBM-Cloud/power-go-client v1.12.0 // indirect
	github.com/IBM/go-sdk-core/v5 v5.21.0 // indirect
	github.com/IBM/ibm-cos-sdk-go v1.12.3 // indirect
	github.com/IBM/keyprotect-go-client v0.12.2 // indirect
	github.com/IBM/networking-go-sdk v0.51.11 // indirect
	github.com/IBM/platform-services-go-sdk v0.86.1 // indirect
	github.com/IBM/vpc-go-sdk v0.71.2 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/PaesslerAG/gval v1.2.2 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/asaskevich/govalidator/v11 v11.0.2-0.20250122183457-e11347878e23 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.5 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.40.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.292.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/efs v1.41.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.32.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.54.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.31.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.96.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.34.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.7 // indirect
	github.com/aws/smithy-go v1.24.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cavaliercoder/go-cpio v0.0.0-20180626203310-925f9528c45e // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/clarketm/json v1.17.1 // indirect
	github.com/cncf/xds/go v0.0.0-20251210132809-ee656c7534f5 // indirect
	github.com/containers/image/v5 v5.31.0 // indirect
	github.com/containers/storage v1.54.0 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/coreos/stream-metadata-go v0.4.10-0.20250806142651-4a7d280a6c7b // indirect
	github.com/coreos/vcontext v0.0.0-20231102161604-685dc7299dc5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/digitalocean/go-libvirt v0.0.0-20240220204746-fcabe97a6eed // indirect
	github.com/diskfs/go-diskfs v1.4.1 // indirect
	github.com/djherbis/times v1.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/elliotwutingfeng/asciiset v0.0.0-20230602022725-51bbb787efab // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.36.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.0 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20210407135951-1de76d718b3f // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/gobuffalo/flect v1.0.3 // indirect
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/gophercloud/gophercloud/v2 v2.8.0 // indirect
	github.com/gophercloud/utils/v2 v2.0.0-20250212084022-725b94822eeb // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.0.12 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.1-0.20220621161143-b0104c826a24 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k-orc/openstack-resource-controller/v2 v2.3.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kdomanski/iso9660 v0.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/metal3-io/baremetal-operator/apis v0.5.1 // indirect
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.5.1 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/microsoft/kiota-abstractions-go v1.9.2 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.3.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nutanix-cloud-native/cluster-api-provider-nutanix v1.7.0 // indirect
	github.com/nutanix-cloud-native/prism-go-client v0.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/openshift/assisted-image-service v0.0.0-20250917153356-4ca9ff81f712 // indirect
	github.com/openshift/assisted-service/api v0.0.0 // indirect
	github.com/openshift/assisted-service/models v0.0.0 // indirect
	github.com/openshift/client-go v0.0.0-20260317180604-743f664b82d1 // indirect
	github.com/openshift/cloud-credential-operator v0.0.0-20250925191450-a3e50aeacc55 // indirect
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20250120030111-6e7aa77e5e14 // indirect
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20240104064818-a336f0b5f7ee // indirect
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20240507153406-abb91c84083d // indirect
	github.com/openshift/custom-resource-status v1.1.3-0.20220503160415-f2fdb4999d87 // indirect
	github.com/openshift/hive/apis v0.0.0-20231220215202-ad99b9e52d27 // indirect
	github.com/openshift/library-go v0.0.0-20251107090138-0de9712313a5 // indirect
	github.com/openshift/machine-api-operator v0.2.1-0.20251110092458-e0af0f3f44b8 // indirect
	github.com/openshift/machine-api-provider-gcp v0.0.1-0.20260113091719-80740861bb2a // indirect
	github.com/openshift/machine-api-provider-ibmcloud v0.0.0-20241121232730-238348d2d06f // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/ovirt/go-ovirt v0.0.0-20220427092237-114c47f2835c // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c // indirect
	github.com/shurcooL/vfsgen v0.0.0-20230704071429-0000e147ea92 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.5 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/thedevsaddam/retry v1.2.1 // indirect
	github.com/thoas/go-funk v0.9.3 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/vmware/govmomi v0.53.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.17.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.39.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.65.0 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/api v0.255.0 // indirect
	google.golang.org/genproto v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/grpc v1.79.3 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.25.7 // indirect
	k8s.io/apiextensions-apiserver v0.35.1 // indirect
	k8s.io/apiserver v0.35.1 // indirect
	k8s.io/cli-runtime v0.35.1 // indirect
	k8s.io/cloud-provider-vsphere v1.32.2 // indirect
	k8s.io/cluster-bootstrap v0.33.3 // indirect
	k8s.io/component-base v0.35.1 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250910181357-589584f1c912 // indirect
	k8s.io/kubectl v0.35.1 // indirect
	k8s.io/utils v0.0.0-20251002143259-bc988d571ff4 // indirect
	libvirt.org/go/libvirtxml v1.10002.0 // indirect
	sigs.k8s.io/cluster-api v1.11.8 // indirect
	sigs.k8s.io/cluster-api-provider-aws/v2 v2.10.1-0.20260226000023-67627b2da0df // indirect
	sigs.k8s.io/cluster-api-provider-gcp v1.10.1-0.20251030112604-5c2bfd5aadb3 // indirect
	sigs.k8s.io/cluster-api-provider-ibmcloud v0.12.2 // indirect
	sigs.k8s.io/cluster-api-provider-openstack v0.13.0 // indirect
	sigs.k8s.io/cluster-api-provider-vsphere v1.14.0 // indirect
	sigs.k8s.io/controller-runtime v0.22.4 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/kustomize/api v0.20.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.20.1 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.2 // indirect
)

// Exclude packages that have weird versioning or are broken
exclude k8s.io/client-go v12.0.0+incompatible

// Dependency pins
replace (
	github.com/imdario/mergo => dario.cat/mergo v0.3.16

	// too-new versions of openshift/api and openshift/client-go don't work
	github.com/openshift/api => github.com/openshift/api v0.0.0-20260825115208-34f4129d8a40
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20260818144019-860af65a7cbc
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.19.3
)

// Installer forks/tags
replace (
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20231128154154-6736c9b9c6c8
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20231128154154-6736c9b9c6c8
)

// They add a require to a module in the same tree with an invalid version, but also add a replace statement to use the local module,
// ie. `github.com/openshift/assisted-service/api` v0.0.0 and `replace github.com/openshift/assisted-service/api => ./api`
// https://github.com/openshift/assisted-service/commit/8c54328a19e2967cc1d111c3a2ad808fc978c546#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6
// The upstream openshift installer adds a replace statement to fix that, but these are not inherited, so we have to add them in the wrapper, too:
// https://github.com/openshift/installer/commit/aed3ddb3bbca729f006efb5499fb80253e492e2b#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6R256
replace (
	github.com/openshift/assisted-service/api => github.com/openshift/assisted-service/api v0.0.0-20260812105124-f8220b56913e
	github.com/openshift/assisted-service/client => github.com/openshift/assisted-service/client v0.0.0-20260812105124-f8220b56913e
	github.com/openshift/assisted-service/models => github.com/openshift/assisted-service/models v0.0.0-20260812105124-f8220b56913e
)

// Temporary vendoring of upstream providers while they work to update to capi v1.11, which contains breaking changes
// See: https://issues.redhat.com/browse/CORS-4262
replace (
	github.com/nutanix-cloud-native/cluster-api-provider-nutanix => github.com/nutanix-cloud-native/cluster-api-provider-nutanix v1.7.2-0.20251007022949-442bc2ebe286
	sigs.k8s.io/cluster-api-provider-azure => github.com/mboersma/cluster-api-provider-azure v0.3.1-0.20251030205607-3161b9cc8d3e
)
