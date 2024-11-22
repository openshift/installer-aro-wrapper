module github.com/openshift/installer-aro-wrapper

go 1.22.0

toolchain go1.22.4

require (
	github.com/AlekSi/gocov-xml v0.0.0-20190121064608-3a14fb1c4737
	github.com/Azure/azure-sdk-for-go v63.1.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.1
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.13
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/alvaroloes/enumer v1.1.2
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/axw/gocov v1.0.0
	github.com/codahale/etm v0.0.0-20141003032925-c00c9e6fb4c9
	github.com/coreos/ignition/v2 v2.20.0
	github.com/coreos/stream-metadata-go v0.4.3
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-test/deep v1.1.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang/mock v1.7.0-rc.1
	github.com/golangci/golangci-lint v1.58.2
	github.com/google/go-cmp v0.6.0
	github.com/jewzaam/go-cosmosdb v0.0.0-20220315232836-282b67c5b234
	github.com/jongio/azidext/go/azidext v0.5.0
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/onsi/gomega v1.34.0
	github.com/openshift/api v0.0.0-20241001152557-e415140e5d5f
	github.com/openshift/client-go v0.0.1
	github.com/openshift/machine-config-operator v3.11.0+incompatible
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	golang.org/x/crypto v0.28.0
	golang.org/x/tools v0.22.0
	gotest.tools/gotestsum v1.6.4
	k8s.io/api v0.31.1
	k8s.io/apimachinery v0.31.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/code-generator v0.31.1
	sigs.k8s.io/controller-tools v0.12.0
)

require (
	cloud.google.com/go/auth v0.7.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.3 // indirect
	cloud.google.com/go/monitoring v1.20.1 // indirect
	github.com/Antonboom/testifylint v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3 v3.0.0-beta.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Crocmagnon/fatcontext v0.2.2 // indirect
	github.com/GaijinEntertainment/go-exhaustruct/v3 v3.2.0 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/alecthomas/go-check-sumtype v0.1.4 // indirect
	github.com/bombsimon/wsl/v4 v4.2.1 // indirect
	github.com/catenacyber/perfsprint v0.7.1 // indirect
	github.com/cavaliercoder/go-cpio v0.0.0-20180626203310-925f9528c45e // indirect
	github.com/ccojocar/zxcvbn-go v1.0.2 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/ckaznocha/intrange v0.1.2 // indirect
	github.com/digitalocean/go-libvirt v0.0.0-20240220204746-fcabe97a6eed // indirect
	github.com/diskfs/go-diskfs v1.4.0 // indirect
	github.com/dougm/pretty v0.0.0-20171025230240-2ee9d7453c02 // indirect
	github.com/elliotwutingfeng/asciiset v0.0.0-20230602022725-51bbb787efab // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/ghostiam/protogetter v0.3.6 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/gofrs/uuid/v5 v5.2.0 // indirect
	github.com/golangci/modinfo v0.3.4 // indirect
	github.com/golangci/plugin-module-register v0.1.1 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/gophercloud/gophercloud/v2 v2.0.0 // indirect
	github.com/gophercloud/utils/v2 v2.0.0-20240701101423-2401526caee5 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/terraform-exec v0.17.3 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/jjti/go-spancheck v0.6.1 // indirect
	github.com/karamaru-alpha/copyloopvar v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lasiar/canonicalheader v1.1.1 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/macabu/inamedparam v0.1.3 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/spdystream v0.4.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/nutanix-cloud-native/cluster-api-provider-nutanix v1.3.3-0.20240416171357-98239ba02cb2 // indirect
	github.com/openshift/assisted-image-service v0.0.0-20240607085136-02df2e56dde6 // indirect
	github.com/openshift/library-go v0.0.0-20240919205913-c96b82b3762b // indirect
	github.com/openshift/machine-api-operator v0.2.1-0.20240722145313-3a817c78946a // indirect
	github.com/openshift/machine-api-provider-gcp v0.0.1-0.20231014045125-6096cc86f3ba // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/ppc64le-cloud/powervs-utils v0.0.0-20240610070307-1c0d75a5c247 // indirect
	github.com/quasilyte/go-ruleguard/dsl v0.3.22 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/zclconf/go-cty v1.11.0 // indirect
	go-simpler.org/musttag v0.12.2 // indirect
	go-simpler.org/sloglint v0.6.0 // indirect
	go.starlark.net v0.0.0-20230525235612-a134d8f9ddca // indirect
	go.uber.org/automaxprocs v1.5.3 // indirect
	gopkg.in/djherbis/times.v1 v1.3.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	k8s.io/cli-runtime v0.31.1 // indirect
	k8s.io/cluster-bootstrap v0.30.3 // indirect
	k8s.io/gengo/v2 v2.0.0-20240228010128-51d4e06bde70 // indirect
	k8s.io/kubectl v0.31.1 // indirect
	libvirt.org/go/libvirtxml v1.10002.0 // indirect
	sigs.k8s.io/cluster-api-provider-gcp v1.7.1-0.20240724153512-c3b8b533143c // indirect
	sigs.k8s.io/cluster-api-provider-ibmcloud v0.7.0 // indirect
	sigs.k8s.io/cluster-api-provider-openstack v0.10.3 // indirect
	sigs.k8s.io/cluster-api-provider-vsphere v1.11.1
	sigs.k8s.io/kustomize/api v0.17.3 // indirect
	sigs.k8s.io/kustomize/kyaml v0.17.2 // indirect
)

require (
	4d63.com/gocheckcompilerdirectives v1.2.1 // indirect
	4d63.com/gochecknoglobals v0.2.1 // indirect
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	cloud.google.com/go/iam v1.1.11 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	github.com/4meepo/tagalign v1.3.4 // indirect
	github.com/Abirdcfly/dupword v0.0.14 // indirect
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/Antonboom/errname v0.1.13 // indirect
	github.com/Antonboom/nilnil v0.1.9 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.9.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.0.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/Djarvur/go-err113 v0.1.0 // indirect
	github.com/IBM-Cloud/bluemix-go v0.0.0-20231026045515-5c31bb9104f7 // indirect
	github.com/IBM-Cloud/power-go-client v1.8.3 // indirect
	github.com/IBM/go-sdk-core/v5 v5.18.1 // indirect
	github.com/IBM/keyprotect-go-client v0.12.2 // indirect
	github.com/IBM/networking-go-sdk v0.45.0 // indirect
	github.com/IBM/platform-services-go-sdk v0.71.0 // indirect
	github.com/IBM/vpc-go-sdk v0.61.0 // indirect
	github.com/OpenPeeDeeP/depguard/v2 v2.2.0 // indirect
	github.com/PaesslerAG/gval v1.2.2 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/alexkohler/nakedret/v2 v2.0.4 // indirect
	github.com/alexkohler/prealloc v1.0.0 // indirect
	github.com/alingse/asasalint v0.0.11 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/ashanbrown/forbidigo v1.6.0 // indirect
	github.com/ashanbrown/makezero v1.1.1 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bkielbasa/cyclop v1.2.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/blizzy78/varnamelen v0.8.0 // indirect
	github.com/breml/bidichk v0.2.7 // indirect
	github.com/breml/errchkjson v0.3.6 // indirect
	github.com/butuzov/ireturn v0.3.0 // indirect
	github.com/butuzov/mirror v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charithe/durationcheck v0.0.10 // indirect
	github.com/chavacava/garif v0.1.0 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/clarketm/json v1.17.1 // indirect
	github.com/containers/image/v5 v5.31.0 // indirect
	github.com/containers/storage v1.54.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/coreos/vcontext v0.0.0-20231102161604-685dc7299dc5 // indirect
	github.com/curioswitch/go-reassign v0.2.0 // indirect
	github.com/daixiang0/gci v0.13.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denis-tingaikin/go-header v0.5.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/ettle/strcase v0.2.0 // indirect
	github.com/evanphx/json-patch/v5 v5.9.0 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/firefart/nonamedreturns v1.0.5 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-critic/go-critic v0.11.3 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.2.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.2 // indirect
	github.com/gobuffalo/flect v1.0.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a // indirect
	github.com/golangci/gofmt v0.0.0-20231018234816-f50ced29576e // indirect
	github.com/golangci/misspell v0.5.1 // indirect
	github.com/golangci/revgrep v0.5.3 // indirect
	github.com/golangci/unconvert v0.0.0-20240309020433-c5143eacb3ed // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gordonklaus/ineffassign v0.1.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.1.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.1 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/imdario/mergo v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jgautheron/goconst v1.7.1 // indirect
	github.com/jingyugao/rowserrcheck v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kdomanski/iso9660 v0.4.0 // indirect
	github.com/kisielk/errcheck v1.7.0 // indirect
	github.com/kkHAIKE/contextcheck v1.1.5 // indirect
	github.com/kulti/thelper v0.6.3 // indirect
	github.com/kunwardeep/paralleltest v1.0.10 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/kyoh86/exportloopref v0.1.11 // indirect
	github.com/ldez/gomoddirectives v0.2.4 // indirect
	github.com/ldez/tagliatelle v0.5.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/leonklingele/grouper v1.1.2 // indirect
	github.com/lufeee/execinquery v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/maratori/testableexamples v1.0.0 // indirect
	github.com/maratori/testpackage v1.1.1 // indirect
	github.com/matoous/godox v0.0.0-20230222163458-006bad1f9d26 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/metal3-io/baremetal-operator/apis v0.4.0 // indirect
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.4.0 // indirect
	github.com/mgechev/revive v1.3.7 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/microsoft/kiota-abstractions-go v1.5.0 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/sys/mountinfo v0.7.1 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/moricho/tparallel v0.3.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nishanths/exhaustive v0.12.0 // indirect
	github.com/nishanths/predeclared v0.2.2 // indirect
	github.com/nunnatsa/ginkgolinter v0.16.2 // indirect
	// This module had a breaking change. Didn't upgrade it. (4.14)
	// https://github.com/nutanix-cloud-native/prism-go-client/pull/30
	github.com/nutanix-cloud-native/prism-go-client v0.3.4 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/openshift/assisted-service/api v0.0.0 // indirect
	github.com/openshift/assisted-service/models v0.0.0 // indirect
	github.com/openshift/cloud-credential-operator v0.0.0-20240422222427-55199c9b5870 // indirect
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20231103125354-2adf4799b2b9 // indirect
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20230308152226-83c0473d4429 // indirect
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20220323121149-e3f2850dd519 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/openshift/hive/apis v0.0.0-20231114185751-51a8b107360f // indirect
	github.com/openshift/machine-api-provider-ibmcloud v0.0.0-20231207164151-6b0b8ea7b16d // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/ovirt/go-ovirt v0.0.0-20220427092237-114c47f2835c // indirect
	github.com/pascaldekloe/name v0.0.0-20180628100202-0fd16699aae1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polyfloyd/go-errorlint v1.5.1 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quasilyte/go-ruleguard v0.4.2 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ryancurrah/gomodguard v1.3.2 // indirect
	github.com/ryanrolds/sqlclosecheck v0.5.1 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.0.7 // indirect
	github.com/sashamelentyev/interfacebloat v1.1.0 // indirect
	github.com/sashamelentyev/usestdlibvars v1.25.0 // indirect
	github.com/securego/gosec/v2 v2.20.0 // indirect
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sivchari/containedctx v1.0.3 // indirect
	github.com/sivchari/tenv v1.7.1 // indirect
	github.com/sonatard/noctx v0.0.2 // indirect
	github.com/sourcegraph/go-diff v0.7.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stbenjam/no-sprintf-host-port v0.1.1 // indirect
	github.com/std-uritemplate/std-uritemplate/go v0.0.46 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c // indirect
	github.com/tdakkota/asciicheck v0.2.0 // indirect
	github.com/tetafro/godot v1.4.16 // indirect
	github.com/thedevsaddam/retry v0.0.0-20200324223450-9769a859cc6d // indirect
	github.com/timakin/bodyclose v0.0.0-20230421092635-574207250966 // indirect
	github.com/timonwong/loggercheck v0.9.4 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.8.3 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.1 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/ultraware/funlen v0.1.0 // indirect
	github.com/ultraware/whitespace v0.1.1 // indirect
	github.com/uudashr/gocognit v1.1.2 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/xen0n/gosmopolitan v1.2.2 // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.3.0 // indirect
	github.com/ykadowak/zerologlint v0.1.5 // indirect
	gitlab.com/bosi/decorder v0.4.2 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/exp/typeparams v0.0.0-20240314144324-c7f7c6466f7f // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/api v0.189.0 // indirect
	google.golang.org/genproto v0.0.0-20240722135656-d784300faade // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240722135656-d784300faade // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.25.5 // indirect
	honnef.co/go/tools v0.4.7 // indirect
	k8s.io/apiextensions-apiserver v0.31.0 // indirect
	k8s.io/cloud-provider-vsphere v1.31.0 // indirect
	k8s.io/component-base v0.31.1 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240921022957-49e7df575cb6 // indirect
	mvdan.cc/gofumpt v0.6.0 // indirect
	mvdan.cc/unparam v0.0.0-20240427195214-063aff900ca1 // indirect
	sigs.k8s.io/cluster-api v1.8.4 // indirect
	sigs.k8s.io/cluster-api-provider-aws/v2 v2.6.1-0.20240910142813-3f3ce5610204 // indirect
	sigs.k8s.io/cluster-api-provider-azure v1.15.1-0.20240617212811-a52056dfb88c // indirect
	sigs.k8s.io/controller-runtime v0.19.0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0
)

exclude (
	// we're not yet compatible with new ones
	github.com/Azure/azure-sdk-for-go v48.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v55.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v63.2.0+incompatible
	github.com/Azure/azure-sdk-for-go v63.3.0+incompatible
	github.com/Azure/azure-sdk-for-go v63.4.0+incompatible
	github.com/Azure/azure-sdk-for-go v64.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v64.1.0+incompatible
	github.com/Azure/azure-sdk-for-go v64.2.0+incompatible
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v66.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v67.0.0+incompatible
	github.com/Azure/azure-sdk-for-go v67.1.0+incompatible
	github.com/Azure/azure-sdk-for-go v67.2.0+incompatible
	github.com/Azure/azure-sdk-for-go v67.3.0+incompatible
	github.com/Azure/azure-sdk-for-go v67.4.0+incompatible
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/containerd/containerd v1.2.7
	// exclude github.com/containerd/containerd < 1.6.1, 1.5.10, 1.14.12 https://nvd.nist.gov/vuln/detail/CVE-2022-23648
	github.com/containerd/containerd v1.2.10
	github.com/containerd/containerd v1.3.0
	github.com/containerd/containerd v1.3.2
	github.com/containerd/containerd v1.4.1
	github.com/containerd/containerd v1.4.3
	github.com/containerd/containerd v1.4.4
	github.com/containerd/containerd v1.4.9
	github.com/containerd/containerd v1.5.0-beta.1
	github.com/containerd/containerd v1.5.0-beta.3
	github.com/containerd/containerd v1.5.0-beta.4
	github.com/containerd/containerd v1.5.0-rc.0
	github.com/containerd/containerd v1.5.1
	github.com/containerd/containerd v1.5.2
	github.com/containerd/containerd v1.5.7
	github.com/containerd/containerd v1.5.9
	// force use of go.etcd.io/bbolt
	github.com/coreos/bbolt v1.3.0
	github.com/coreos/bbolt v1.3.2
	github.com/coreos/bbolt v1.3.3
	github.com/etcd-io/bbolt v1.3.3
	github.com/etcd-io/bbolt v1.3.6
	// exclude github.com/golang/protobuf < 1.3.2 https://nvd.nist.gov/vuln/detail/CVE-2021-3121
	github.com/gogo/protobuf v1.0.0
	github.com/gogo/protobuf v1.1.1
	github.com/gogo/protobuf v1.2.0
	github.com/gogo/protobuf v1.2.1
	github.com/gogo/protobuf v1.3.0
	github.com/gogo/protobuf v1.3.1
	// force use of golang.org/x/lint
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7
	// exclude github.com/golang/protobuf < 1.3.2 https://nvd.nist.gov/vuln/detail/CVE-2021-3121
	github.com/golang/protobuf v1.0.0
	github.com/golang/protobuf v1.1.1
	github.com/golang/protobuf v1.2.0
	github.com/golang/protobuf v1.2.1
	github.com/golang/protobuf v1.3.0
	github.com/golang/protobuf v1.3.1
	// exclude github.com/hashicorp/vault < v1.5.1: https://nvd.nist.gov/vuln/detail/CVE-2020-16251
	github.com/hashicorp/vault v0.10.4
	github.com/hashicorp/vault v0.11.0-beta1
	github.com/hashicorp/vault v0.11.0
	github.com/hashicorp/vault v0.11.1
	github.com/hashicorp/vault v0.11.2
	github.com/hashicorp/vault v0.11.3
	github.com/hashicorp/vault v0.11.4
	github.com/hashicorp/vault v0.11.5
	github.com/hashicorp/vault v0.11.6
	github.com/hashicorp/vault v0.11.7
	github.com/hashicorp/vault v0.11.8
	github.com/hashicorp/vault v1.0.0-beta1
	github.com/hashicorp/vault v1.0.0-beta2
	github.com/hashicorp/vault v1.0.0-rc1
	github.com/hashicorp/vault v1.0.0
	github.com/hashicorp/vault v1.0.1
	github.com/hashicorp/vault v1.0.2
	github.com/hashicorp/vault v1.0.3
	github.com/hashicorp/vault v1.1.0-beta1
	github.com/hashicorp/vault v1.1.0-beta2
	github.com/hashicorp/vault v1.1.0
	github.com/hashicorp/vault v1.1.1
	github.com/hashicorp/vault v1.1.2
	github.com/hashicorp/vault v1.1.3
	github.com/hashicorp/vault v1.1.4
	github.com/hashicorp/vault v1.1.5
	github.com/hashicorp/vault v1.2.0-beta1
	github.com/hashicorp/vault v1.2.0-beta2
	github.com/hashicorp/vault v1.2.0-rc1
	github.com/hashicorp/vault v1.2.0
	github.com/hashicorp/vault v1.2.1
	github.com/hashicorp/vault v1.2.2
	github.com/hashicorp/vault v1.2.3
	github.com/hashicorp/vault v1.2.4
	github.com/hashicorp/vault v1.3.0
	github.com/hashicorp/vault v1.3.1
	github.com/hashicorp/vault v1.3.2
	github.com/hashicorp/vault v1.3.3
	github.com/hashicorp/vault v1.3.4
	github.com/hashicorp/vault v1.3.5
	github.com/hashicorp/vault v1.3.6
	github.com/hashicorp/vault v1.3.7
	github.com/hashicorp/vault v1.4.0-beta1
	github.com/hashicorp/vault v1.4.0-beta2
	github.com/hashicorp/vault v1.4.0-beta3
	github.com/hashicorp/vault v1.4.0-rc1
	github.com/hashicorp/vault v1.4.0
	github.com/hashicorp/vault v1.4.1
	github.com/hashicorp/vault v1.4.2
	github.com/hashicorp/vault v1.4.3
	github.com/hashicorp/vault v1.5.0-beta1
	github.com/hashicorp/vault v1.5.0-beta2
	github.com/hashicorp/vault v1.5.0-rc1
	github.com/hashicorp/vault v1.5.0
	// force use of github.com/envoyproxy/protoc-gen-validate
	github.com/lyft/protoc-gen-validate v0.0.13
	// busted install path
	github.com/mikefarah/yaml/v2 v2.4.0
	// https://www.whitesourcesoftware.com/vulnerability-database/WS-2018-0594
	github.com/satori/go.uuid v0.0.0
	github.com/satori/uuid v0.0.0
	// force use of cloud.google.com/go
	google.golang.org/cloud v0.0.0-20151119220103-975617b05ea8
)

replace (
	bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d // 404 on bitbucket.org/ww/goautoneg
	github.com/Unknwon/com => github.com/unknwon/com v1.0.1
	github.com/clarketm/json => github.com/clarketm/json v1.15.7 // Later versions not compatible with Go 1.16
	github.com/cockroachdb/sentry-go => github.com/getsentry/sentry-go v0.11.0
	github.com/docker/spdystream => github.com/docker/spdystream v0.1.0
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.20.4
	// Replace old GoGo Protobuf versions https://nvd.nist.gov/vuln/detail/CVE-2021-3121
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/mrnold/go-libnbd => github.com/mrnold/go-libnbd v1.4.1-cdi // v1.10.0 uses an invalid module path
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v1.21.0 // earlier versions refer to dead packages
	// https://www.whitesourcesoftware.com/vulnerability-database/WS-2018-0594
	github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/satori/uuid => github.com/satori/uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace
	github.com/spf13/viper => github.com/spf13/viper v1.7.1
	github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20200630224953-76d1fb4e5699
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.40.1-0.20200707062554-97ea089cc12a
	github.com/terraform-providers/terraform-provider-ignition/v2 => github.com/community-terraform-providers/terraform-provider-ignition/v2 v2.1.0
	k8s.io/api => k8s.io/api v0.31.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.31.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.31.1
	k8s.io/client-go => k8s.io/client-go v0.31.1
	k8s.io/code-generator => k8s.io/code-generator v0.31.1
	k8s.io/component-base => k8s.io/component-base v0.31.1
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.5.0
)

// Installer dependencies. Some of them are being used directly in the RP.
replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.15.0
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/BurntSushi/toml => github.com/BurntSushi/toml v0.3.1
	github.com/IBM-Cloud/terraform-provider-ibm => github.com/openshift/terraform-provider-ibm v1.26.2-openshift-2
	github.com/c-bata/go-prompt => github.com/c-bata/go-prompt v0.2.5
	github.com/circonus-labs/circonusllhist => github.com/openhistogram/circonusllhist v0.3.0
	github.com/cockroachdb/errors => github.com/cockroachdb/errors v1.8.5
	github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/containernetworking/plugins => github.com/containernetworking/plugins v1.0.0
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.6
	github.com/coreos/fcct => github.com/coreos/butane v0.13.1
	github.com/coreos/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.48.1
	github.com/coreos/stream-metadata-go => github.com/coreos/stream-metadata-go v0.1.8
	github.com/cortexproject/cortex => github.com/cortexproject/cortex v1.10.0
	github.com/deislabs/oras => github.com/oras-project/oras v0.12.0
	github.com/etcd-io/bbolt => go.etcd.io/bbolt v1.3.6
	github.com/go-check/check => gopkg.in/check.v1 v0.0.0-20201130134442-10cb98267c6c
	github.com/golang/lint => golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	github.com/google/tcpproxy => inet.af/tcpproxy v0.0.0-20210824174053-2e577fef49e2
	github.com/googleapis/gnostic => github.com/google/gnostic v0.5.5
	github.com/h2non/filetype => github.com/h2non/filetype v1.1.1
	github.com/hashicorp/vault => github.com/hasicorp/vault v1.8.7
	github.com/imdario/mergo => dario.cat/mergo v0.3.16
	github.com/influxdata/flux => github.com/influxdata/flux v0.132.0
	github.com/knq/sysutil => github.com/chromedp/sysutil v1.0.0
	github.com/kshvakov/clickhouse => github.com/ClickHouse/clickhouse-go v1.4.9
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20240904104510-62a7d56e2078 // Use OpenShift fork
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20240904104510-62a7d56e2078 // Use OpenShift fork
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20240904104510-62a7d56e2078 // Use OpenShift fork
	github.com/metal3-io/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20240326141158-73a72cc82122 // Pin OpenShift fork
	github.com/metal3-io/cluster-api-provider-metal3 => github.com/openshift/cluster-api-provider-metal3 v0.0.0-20231025134231-298f47e3f7e4
	github.com/metal3-io/cluster-api-provider-metal3/api => github.com/openshift/cluster-api-provider-metal3/api v0.0.0-20231025134231-298f47e3f7e4
	github.com/mholt/certmagic => github.com/caddyserver/certmagic v0.15.0
	github.com/openshift/api => github.com/openshift/api v0.0.0-20241007111039-82e082220d91
	github.com/openshift/assisted-service => github.com/openshift/assisted-service v1.0.10-0.20241011223353-213c57977cc8
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20241001162912-da6d55e4611f
	github.com/openshift/cloud-credential-operator => github.com/openshift/cloud-credential-operator v0.0.0-20240422222427-55199c9b5870
	github.com/openshift/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20240326141158-73a72cc82122
	github.com/openshift/cluster-api-provider-gcp => github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20211123160814-0d569513f9fa
	github.com/openshift/cluster-api-provider-kubevirt => github.com/openshift/cluster-api-provider-kubevirt v0.0.0-20210719100556-9b8bc3666720
	github.com/openshift/cluster-api-provider-libvirt => github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20240104073218-1e096cdf1bfd
	github.com/openshift/cluster-api-provider-ovirt => github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20230728074040-5d708631fca3
	github.com/openshift/console-operator => github.com/openshift/console-operator v0.0.0-20220902005223-378917170147
	github.com/openshift/installer => github.com/openshift/installer v1.4.17
	github.com/openshift/library-go => github.com/openshift/library-go v0.0.0-20240919205913-c96b82b3762b
	github.com/openshift/machine-api-operator => github.com/openshift/machine-api-operator v0.2.1-0.20240930121047-57b7917e6140
	github.com/openshift/machine-api-provider-ibmcloud => github.com/openshift/machine-api-provider-ibmcloud v0.0.0-20231207164151-6b0b8ea7b16d
	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20240927143910-dc45d313d876
	github.com/oras-project/oras-go => oras.land/oras-go v0.4.0
	github.com/ovirt/go-ovirt => github.com/ovirt/go-ovirt v0.0.0-20210112072624-e4d3b104de71
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20210421143221-52df5ef7a3be
	github.com/terraform-providers/terraform-provider-azuread => github.com/hashicorp/terraform-provider-azuread v1.6.0
	github.com/thanos-io/thanos => github.com/thanos-io/thanos v0.23.0
	github.com/uber-go/atomic => go.uber.org/atomic v1.9.0
	github.com/uber/athenadriver => github.com/uber/athenadriver v1.1.10
	github.com/vmware/govmomi => github.com/vmware/govmomi v0.39.0
	github.com/willf/bitset => github.com/bits-and-blooms/bitset v1.2.1
	google.golang.org/cloud => cloud.google.com/go v0.97.0
	// CVE-2023-44487 gRPC-Go HTTP/2 Rapid Reset vulnerability https://github.com/grpc/grpc-go/pull/6703
	google.golang.org/grpc => google.golang.org/grpc v1.56.3
	k8s.io/cloud-provider-vsphere => github.com/openshift/cloud-provider-vsphere v1.19.1-0.20240626105621-6464d0bb4928
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.70.1
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.9.7
	mvdan.cc/unparam => mvdan.cc/unparam v0.0.0-20211002133954-f839ab2b2b11
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.8.4
	sigs.k8s.io/cluster-api-provider-aws/v2 => sigs.k8s.io/cluster-api-provider-aws/v2 v2.6.1-0.20241026111253-5b4f7c1acb52
	sigs.k8s.io/cluster-api-provider-azure => sigs.k8s.io/cluster-api-provider-azure v1.15.1-0.20240617212811-a52056dfb88c
	sigs.k8s.io/cluster-api-provider-ibmcloud => sigs.k8s.io/cluster-api-provider-ibmcloud v0.9.0-beta.0.0.20241029051454-9b0770491a76
	sigs.k8s.io/kustomize/api => sigs.k8s.io/kustomize/api v0.17.3
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.17.2
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
	vbom.ml/util => github.com/fvbommel/util v0.0.3
)

// They add a require to a module in the same tree with an invalid version, but also add a replace statement to use the local module,
// ie. `github.com/openshift/assisted-service/api` v0.0.0 and `replace github.com/openshift/assisted-service/api => ./api`
// https://github.com/openshift/assisted-service/commit/8c54328a19e2967cc1d111c3a2ad808fc978c546#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6
// The upstream openshift installer adds a replace statement to fix that, but these are not inherited, so we have to add them in the wrapper, too:
// https://github.com/openshift/installer/commit/aed3ddb3bbca729f006efb5499fb80253e492e2b#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6R256
replace (
	github.com/openshift/assisted-service/api => github.com/openshift/assisted-service/api v0.0.0-20230831114549-1922eda29cf8
	github.com/openshift/assisted-service/client => github.com/openshift/assisted-service/client v0.0.0-20230831114549-1922eda29cf8
	github.com/openshift/assisted-service/models => github.com/openshift/assisted-service/models v0.0.0-20230831114549-1922eda29cf8
	github.com/openshift/hive => github.com/openshift/hive v1.1.17-0.20220719141355-c63c9b0281d8
	github.com/openshift/hive/apis => github.com/openshift/hive/apis v0.0.0-20220719141355-c63c9b0281d8
)
