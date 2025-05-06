package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

var expectedIgnitionFileContents = map[string]string{
	"/etc/NetworkManager/dispatcher.d/99-dnsmasq-restart": `
#!/bin/sh
# This is a NetworkManager dispatcher script to restart dnsmasq
# in the event of a network interface change (e. g. host servicing event https://learn.microsoft.com/en-us/azure/developer/intro/hosting-apps-on-azure)
# this will restart dnsmasq, reapplying our /etc/resolv.conf file and overwriting any modifications made by NetworkManager

interface=$1
action=$2

log() {
    logger -i "$0" -t '99-DNSMASQ-RESTART SCRIPT' "$@"
}

# log dns configuration information relevant to SRE while troubleshooting
# The line break used here is important for formatting
log_dns_files() {
    log "/etc/resolv.conf contents

    $(cat /etc/resolv.conf)"

    log "$(echo -n \"/etc/resolv.conf file metadata: \") $(ls -lZ /etc/resolv.conf)"

    log "/etc/resolv.conf.dnsmasq contents

    $(cat /etc/resolv.conf.dnsmasq)"

    log "$(echo -n "/etc/resolv.conf.dnsmasq file metadata: ") $(ls -lZ /etc/resolv.conf.dnsmasq)"
}

if [[ $interface == eth* && $action == "up" ]] || [[ $interface == eth* && $action == "down" ]] || [[ $interface == enP* && $action == "up" ]] || [[ $interface == enP* && $action == "down" ]]; then
    log "$action happened on $interface, connection state is now $CONNECTIVITY_STATE"
    log "Pre dnsmasq restart file information"
    log_dns_files
    log "restarting dnsmasq now"
    if systemctl try-restart dnsmasq --wait; then
        log "dnsmasq successfully restarted"
        log "Post dnsmasq restart file information"
        log_dns_files
    else
        log "failed to restart dnsmasq"
    fi
fi

exit 0
`,
	"/etc/dnsmasq.conf": `
resolv-file=/etc/resolv.conf.dnsmasq
dns-forward-max=10000
address=/api.test-cluster.test.example.com/203.0.113.1
address=/api-int.test-cluster.test.example.com/203.0.113.1
address=/.apps.test-cluster.test.example.com/192.0.2.1
address=/gateway.mock1.example.com/203.0.113.2
address=/gateway.mock2.example.com/203.0.113.2
user=dnsmasq
group=dnsmasq
no-hosts
cache-size=0
`,
	"/etc/fluentbit/journal.conf": `[INPUT]
	Name systemd
	Tag journald
	DB /var/lib/fluent/journald

[FILTER]
	Name modify
	Match journald
	Remove_wildcard _
	Remove TIMESTAMP
	Remove SYSLOG_FACILITY

[OUTPUT]
	Name forward
	Port 24224
`,
	"/etc/hosts.d/aro.conf": `203.0.113.1	api.test-cluster.test.example.com api-int.test-cluster.test.example.com
203.0.113.2	gateway.mock1.example.com gateway.mock2.example.com
`,
	"/etc/mdsd.d/mdsd.env": `MONITORING_GCS_ENVIRONMENT=test-logging-environment
MONITORING_GCS_ACCOUNT=test-logging-account
MONITORING_GCS_REGION=centralus
MONITORING_GCS_CERT_CERTFILE=/etc/mdsd.d/secret/mdsdcert.pem
MONITORING_GCS_CERT_KEYFILE=/etc/mdsd.d/secret/mdsdcert.pem
MONITORING_GCS_NAMESPACE=test-logging-namespace
MONITORING_CONFIG_VERSION=42
MONITORING_USE_GENEVA_CONFIG_SERVICE=true
MONITORING_TENANT=centralus
MONITORING_ROLE=cluster
MONITORING_ROLE_INSTANCE=bootstrap
RESOURCE_ID=test-cluster-resource-id
SUBSCRIPTION_ID=test-subscription
RESOURCE_GROUP=test-resource-group
RESOURCE_NAME=test-logging-resource
`,
	"/etc/mdsd.d/secret/mdsdcert.pem": `# This is not a real private key
# This is not a real certificate
`,
	"/etc/sysconfig/fluentbit": `FLUENTIMAGE=registry.example.com/fluentbit:latest
`,
	"/etc/sysconfig/mdsd": `MDSDIMAGE=registry.example.com/mdsd:latest
`,
	"/opt/openshift/manifests/aro-imageregistry.yaml": `apiVersion: imageregistry.operator.openshift.io/v1
kind: Config
metadata:
  finalizers:
  - imageregistry.operator.openshift.io/finalizer
  name: cluster
spec:
  httpSecret: "test"
  managementState: Managed
  replicas: 2
  storage:
    azure:
      accountName: "test-image-registry-storage-acct"
      container: "image-registry"
      cloudName: "AzurePublicCloud"
    managementState: Unmanaged
`,
	"/opt/openshift/manifests/aro-ingress-service.yaml": `apiVersion: v1
kind: Service
metadata:
  name: router-default
  namespace: openshift-ingress
  annotations:
    service.beta.kubernetes.io/azure-load-balancer-internal: ""
    service.beta.kubernetes.io/azure-load-balancer-ipv4: "192.0.2.1"
  labels:
    app: router
    ingresscontroller.operator.openshift.io/owning-ingresscontroller: default
spec:
  externalTrafficPolicy: Local
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  - name: https
    port: 443
    protocol: TCP
    targetPort: https
  selector:
    ingresscontroller.operator.openshift.io/deployment-ingresscontroller: default
  type: LoadBalancer
`,
	"/etc/NetworkManager/dispatcher.d/30-eth0-mtu-3900": `#!/bin/bash

if [ "$1" == "eth0" ] && [ "$2" == "up" ]; then
    ip link set $1 mtu 3900
fi`,
	"/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-dns.yaml": `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-aro-dns
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      config:
        replace:
          verification: {}
      proxy: {}
      security:
        tls: {}
      timeouts: {}
      version: 3.2.0
    passwd: {}
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,CnJlc29sdi1maWxlPS9ldGMvcmVzb2x2LmNvbmYuZG5zbWFzcQpkbnMtZm9yd2FyZC1tYXg9MTAwMDAKYWRkcmVzcz0vYXBpLnRlc3QtY2x1c3Rlci50ZXN0LmV4YW1wbGUuY29tLzIwMy4wLjExMy4xCmFkZHJlc3M9L2FwaS1pbnQudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20vMjAzLjAuMTEzLjEKYWRkcmVzcz0vLmFwcHMudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20vMTkyLjAuMi4xCmFkZHJlc3M9L2dhdGV3YXkubW9jazEuZXhhbXBsZS5jb20vMjAzLjAuMTEzLjIKYWRkcmVzcz0vZ2F0ZXdheS5tb2NrMi5leGFtcGxlLmNvbS8yMDMuMC4xMTMuMgp1c2VyPWRuc21hc3EKZ3JvdXA9ZG5zbWFzcQpuby1ob3N0cwpjYWNoZS1zaXplPTAK
          verification: {}
        group: {}
        mode: 420
        overwrite: true
        path: /etc/dnsmasq.conf
        user:
          name: root
      - contents:
          source: data:text/plain;charset=utf-8;base64,CiMhL2Jpbi9iYXNoCnNldCAtZXVvIHBpcGVmYWlsCgojIFRoaXMgYmFzaCBzY3JpcHQgaXMgYSBwYXJ0IG9mIHRoZSBBUk8gRG5zTWFzcSBjb25maWd1cmF0aW9uCiMgSXQncyBkZXBsb3llZCBhcyBwYXJ0IG9mIHRoZSA5OS1hcm8tZG5zLSogbWFjaGluZSBjb25maWcKIyBTZWUgaHR0cHM6Ly9naXRodWIuY29tL0F6dXJlL0FSTy1SUAoKIyBUaGlzIGZpbGUgY2FuIGJlIHJlcnVuIGFuZCB0aGUgZWZmZWN0IGlzIGlkZW1wb3RlbnQsIG91dHB1dCBtaWdodCBjaGFuZ2UgaWYgdGhlIERIQ1AgY29uZmlndXJhdGlvbiBjaGFuZ2VzCgpOT0RFSVA9JCgvc2Jpbi9pcCAtLWpzb24gcm91dGUgZ2V0IDE2OC42My4xMjkuMTYgfCAvYmluL2pxIC1yICIuW10ucHJlZnNyYyIpCgppZiBbICIkTk9ERUlQIiAhPSAiIiBdOyB0aGVuCiAgICAvYmluL2NwIC1aIC9ldGMvcmVzb2x2LmNvbmYgL2V0Yy9yZXNvbHYuY29uZi5kbnNtYXNxCiAgICBTRUFSQ0hET01BSU49JChhd2sgJy9ec2VhcmNoLyB7IHByaW50ICQyOyB9JyAvZXRjL3Jlc29sdi5jb25mLmRuc21hc3EpCiAgICAvYmluL2NobW9kIDA3NDQgL2V0Yy9yZXNvbHYuY29uZi5kbnNtYXNxCgogICAgY2F0IDw8RU9GIHwgL2Jpbi90ZWUgL2V0Yy9OZXR3b3JrTWFuYWdlci9jb25mLmQvYXJvLWRucy5jb25mCiMgQWRkZWQgYnkgZG5zbWFzcS5zZXJ2aWNlCltnbG9iYWwtZG5zXQpzZWFyY2hlcz0kU0VBUkNIRE9NQUlOCgpbZ2xvYmFsLWRucy1kb21haW4tKl0Kc2VydmVycz0kTk9ERUlQCkVPRgoKICAgICMgbmV0d29yayBtYW5hZ2VyIG1heSBhbHJlYWR5IGJlIHJ1bm5pbmcgYXQgdGhpcyBwb2ludC4KICAgICMgcmVsb2FkIHRvIHVwZGF0ZSAvZXRjL3Jlc29sdi5jb25mIHdpdGggdGhpcyBjb25maWd1cmF0aW9uCiAgICAvdXNyL2Jpbi9ubWNsaSBnZW5lcmFsIHJlbG9hZCBjb25mCiAgICAvdXNyL2Jpbi9ubWNsaSBnZW5lcmFsIHJlbG9hZCBkbnMtcmMKZmkK
          verification: {}
        group: {}
        mode: 484
        overwrite: true
        path: /usr/local/bin/aro-dnsmasq-pre.sh
        user:
          name: root
    systemd:
      units:
      - contents: |2

          [Unit]
          Description=DNS caching server.
          After=network-online.target
          Before=bootkube.service

          [Service]
          # ExecStartPre will create a copy of the customer current resolv.conf file and make it upstream DNS.
          # This file is a product of user DNS settings on the VNET. We will replace this file to point to
          # dnsmasq instance on the node. dnsmasq will inject certain dns records we need and forward rest of the queries to
          # resolv.conf.dnsmasq upstream customer dns.
          ExecStartPre=/bin/bash /usr/local/bin/aro-dnsmasq-pre.sh
          ExecStart=/usr/sbin/dnsmasq -k
          ExecStopPost=/bin/bash -c '/bin/rm /etc/NetworkManager/conf.d/aro-dns.conf && /usr/bin/nmcli general reload conf && /usr/bin/nmcli general reload dns-rc'
          Restart=always
          StandardOutput=journal+console
          StandardError=journal+console

          [Install]
          WantedBy=multi-user.target
        enabled: true
        name: dnsmasq.service
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`,
	"/opt/openshift/manifests/aro-ingress-namespace.yaml": `apiVersion: v1
kind: Namespace
metadata:
  name: openshift-ingress
`,
	"/opt/openshift/manifests/aro-worker-registries.yaml": `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: 90-aro-worker-registries
spec:
  config:
    ignition:
      version: 2.2.0
    storage:
      files:
      - contents:
          source: "data:text/plain,unqualified-search-registries%20%3D%20%5B%22registry.access.redhat.com%22%2C%20%22docker.io%22%5D%0A%0A%5B%5Bregistry%5D%5D%0A%20%20prefix%20%3D%20%22%22%0A%20%20location%20%3D%20%22quay.io%2Fopenshift-release-dev%2Focp-release%22%0A%20%20mirror-by-digest-only%20%3D%20true%0A%0A%20%20%5B%5Bregistry.mirror%5D%5D%0A%20%20%20%20location%20%3D%20%22registry.example.com%2Fopenshift-release-dev%2Focp-release%22%0A%0A%5B%5Bregistry%5D%5D%0A%20%20prefix%20%3D%20%22%22%0A%20%20location%20%3D%20%22quay.io%2Fopenshift-release-dev%2Focp-release-nightly%22%0A%20%20mirror-by-digest-only%20%3D%20true%0A%0A%20%20%5B%5Bregistry.mirror%5D%5D%0A%20%20%20%20location%20%3D%20%22registry.example.com%2Fopenshift-release-dev%2Focp-release-nightly%22%0A%0A%5B%5Bregistry%5D%5D%0A%20%20prefix%20%3D%20%22%22%0A%20%20location%20%3D%20%22quay.io%2Fopenshift-release-dev%2Focp-v4.0-art-dev%22%0A%20%20mirror-by-digest-only%20%3D%20true%0A%0A%20%20%5B%5Bregistry.mirror%5D%5D%0A%20%20%20%20location%20%3D%20%22registry.example.com%2Fopenshift-release-dev%2Focp-v4.0-art-dev%22%0A"
        filesystem: root
        mode: 420
        path: /etc/containers/registries.conf
`,
	"/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-etc-hosts-gateway-domains.yaml": `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-aro-etc-hosts-gateway-domains
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,MjAzLjAuMTEzLjEJYXBpLnRlc3QtY2x1c3Rlci50ZXN0LmV4YW1wbGUuY29tIGFwaS1pbnQudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20KMjAzLjAuMTEzLjIJZ2F0ZXdheS5tb2NrMS5leGFtcGxlLmNvbSBnYXRld2F5Lm1vY2syLmV4YW1wbGUuY29tCg==
        mode: 420
        overwrite: true
        path: /etc/hosts.d/aro.conf
        user:
          name: root
      - contents:
          source: data:text/plain;charset=utf-8;base64,IyEvYmluL2Jhc2gKc2V0IC11byBwaXBlZmFpbAoKdHJhcCAnam9icyAtcCB8IHhhcmdzIGtpbGwgfHwgdHJ1ZTsgd2FpdDsgZXhpdCAwJyBURVJNCgpPUEVOU0hJRlRfTUFSS0VSPSJvcGVuc2hpZnQtYXJvLWV0Y2hvc3RzLXJlc29sdmVyIgpIT1NUU19GSUxFPSIvZXRjL2hvc3RzIgpDT05GSUdfRklMRT0iL2V0Yy9ob3N0cy5kL2Fyby5jb25mIgpURU1QX0ZJTEU9Ii9ldGMvaG9zdHMuZC9hcm8udG1wIgoKIyBNYWtlIGEgdGVtcG9yYXJ5IGZpbGUgd2l0aCB0aGUgb2xkIGhvc3RzIGZpbGUncyBkYXRhLgppZiAhIGNwIC1mICIke0hPU1RTX0ZJTEV9IiAiJHtURU1QX0ZJTEV9IjsgdGhlbgogIGVjaG8gIkZhaWxlZCB0byBwcmVzZXJ2ZSBob3N0cyBmaWxlLiBFeGl0aW5nLiIKICBleGl0IDEKZmkKCmlmICEgc2VkIC0tc2lsZW50ICIvIyAke09QRU5TSElGVF9NQVJLRVJ9L2Q7IHcgJHtURU1QX0ZJTEV9IiAiJHtIT1NUU19GSUxFfSI7IHRoZW4KICAjIE9ubHkgY29udGludWUgcmVidWlsZGluZyB0aGUgaG9zdHMgZW50cmllcyBpZiBpdHMgb3JpZ2luYWwgY29udGVudCBpcyBwcmVzZXJ2ZWQKICBzbGVlcCA2MCAmIHdhaXQKICBjb250aW51ZQpmaQoKd2hpbGUgSUZTPSByZWFkIC1yIGxpbmU7IGRvCiAgICBlY2hvICIke2xpbmV9ICMgJHtPUEVOU0hJRlRfTUFSS0VSfSIgPj4gIiR7VEVNUF9GSUxFfSIKZG9uZSA8ICIke0NPTkZJR19GSUxFfSIKCiMgUmVwbGFjZSAvZXRjL2hvc3RzIHdpdGggb3VyIG1vZGlmaWVkIHZlcnNpb24gaWYgbmVlZGVkCmNtcCAiJHtURU1QX0ZJTEV9IiAiJHtIT1NUU19GSUxFfSIgfHwgY3AgLWYgIiR7VEVNUF9GSUxFfSIgIiR7SE9TVFNfRklMRX0iCiMgVEVNUF9GSUxFIGlzIG5vdCByZW1vdmVkIHRvIGF2b2lkIGZpbGUgY3JlYXRlL2RlbGV0ZSBhbmQgYXR0cmlidXRlcyBjb3B5IGNodXJuCg==
        mode: 484
        overwrite: true
        path: /usr/local/bin/aro-etchosts-resolver.sh
        user:
          name: root
    systemd:
      units:
      - contents: |
          [Unit]
          Description=One shot service that appends static domains to etchosts
          Before=network-online.target

          [Service]
          # ExecStart will copy the hosts defined in /etc/hosts.d/aro.conf to /etc/hosts
          ExecStart=/bin/bash /usr/local/bin/aro-etchosts-resolver.sh

          [Install]
          WantedBy=multi-user.target
        enabled: true
        name: aro-etchosts-resolver.service
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`,
	"/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-dns.yaml": `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-aro-dns
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      config:
        replace:
          verification: {}
      proxy: {}
      security:
        tls: {}
      timeouts: {}
      version: 3.2.0
    passwd: {}
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,CnJlc29sdi1maWxlPS9ldGMvcmVzb2x2LmNvbmYuZG5zbWFzcQpkbnMtZm9yd2FyZC1tYXg9MTAwMDAKYWRkcmVzcz0vYXBpLnRlc3QtY2x1c3Rlci50ZXN0LmV4YW1wbGUuY29tLzIwMy4wLjExMy4xCmFkZHJlc3M9L2FwaS1pbnQudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20vMjAzLjAuMTEzLjEKYWRkcmVzcz0vLmFwcHMudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20vMTkyLjAuMi4xCmFkZHJlc3M9L2dhdGV3YXkubW9jazEuZXhhbXBsZS5jb20vMjAzLjAuMTEzLjIKYWRkcmVzcz0vZ2F0ZXdheS5tb2NrMi5leGFtcGxlLmNvbS8yMDMuMC4xMTMuMgp1c2VyPWRuc21hc3EKZ3JvdXA9ZG5zbWFzcQpuby1ob3N0cwpjYWNoZS1zaXplPTAK
          verification: {}
        group: {}
        mode: 420
        overwrite: true
        path: /etc/dnsmasq.conf
        user:
          name: root
      - contents:
          source: data:text/plain;charset=utf-8;base64,CiMhL2Jpbi9iYXNoCnNldCAtZXVvIHBpcGVmYWlsCgojIFRoaXMgYmFzaCBzY3JpcHQgaXMgYSBwYXJ0IG9mIHRoZSBBUk8gRG5zTWFzcSBjb25maWd1cmF0aW9uCiMgSXQncyBkZXBsb3llZCBhcyBwYXJ0IG9mIHRoZSA5OS1hcm8tZG5zLSogbWFjaGluZSBjb25maWcKIyBTZWUgaHR0cHM6Ly9naXRodWIuY29tL0F6dXJlL0FSTy1SUAoKIyBUaGlzIGZpbGUgY2FuIGJlIHJlcnVuIGFuZCB0aGUgZWZmZWN0IGlzIGlkZW1wb3RlbnQsIG91dHB1dCBtaWdodCBjaGFuZ2UgaWYgdGhlIERIQ1AgY29uZmlndXJhdGlvbiBjaGFuZ2VzCgpOT0RFSVA9JCgvc2Jpbi9pcCAtLWpzb24gcm91dGUgZ2V0IDE2OC42My4xMjkuMTYgfCAvYmluL2pxIC1yICIuW10ucHJlZnNyYyIpCgppZiBbICIkTk9ERUlQIiAhPSAiIiBdOyB0aGVuCiAgICAvYmluL2NwIC1aIC9ldGMvcmVzb2x2LmNvbmYgL2V0Yy9yZXNvbHYuY29uZi5kbnNtYXNxCiAgICBTRUFSQ0hET01BSU49JChhd2sgJy9ec2VhcmNoLyB7IHByaW50ICQyOyB9JyAvZXRjL3Jlc29sdi5jb25mLmRuc21hc3EpCiAgICAvYmluL2NobW9kIDA3NDQgL2V0Yy9yZXNvbHYuY29uZi5kbnNtYXNxCgogICAgY2F0IDw8RU9GIHwgL2Jpbi90ZWUgL2V0Yy9OZXR3b3JrTWFuYWdlci9jb25mLmQvYXJvLWRucy5jb25mCiMgQWRkZWQgYnkgZG5zbWFzcS5zZXJ2aWNlCltnbG9iYWwtZG5zXQpzZWFyY2hlcz0kU0VBUkNIRE9NQUlOCgpbZ2xvYmFsLWRucy1kb21haW4tKl0Kc2VydmVycz0kTk9ERUlQCkVPRgoKICAgICMgbmV0d29yayBtYW5hZ2VyIG1heSBhbHJlYWR5IGJlIHJ1bm5pbmcgYXQgdGhpcyBwb2ludC4KICAgICMgcmVsb2FkIHRvIHVwZGF0ZSAvZXRjL3Jlc29sdi5jb25mIHdpdGggdGhpcyBjb25maWd1cmF0aW9uCiAgICAvdXNyL2Jpbi9ubWNsaSBnZW5lcmFsIHJlbG9hZCBjb25mCiAgICAvdXNyL2Jpbi9ubWNsaSBnZW5lcmFsIHJlbG9hZCBkbnMtcmMKZmkK
          verification: {}
        group: {}
        mode: 484
        overwrite: true
        path: /usr/local/bin/aro-dnsmasq-pre.sh
        user:
          name: root
    systemd:
      units:
      - contents: |2

          [Unit]
          Description=DNS caching server.
          After=network-online.target
          Before=bootkube.service

          [Service]
          # ExecStartPre will create a copy of the customer current resolv.conf file and make it upstream DNS.
          # This file is a product of user DNS settings on the VNET. We will replace this file to point to
          # dnsmasq instance on the node. dnsmasq will inject certain dns records we need and forward rest of the queries to
          # resolv.conf.dnsmasq upstream customer dns.
          ExecStartPre=/bin/bash /usr/local/bin/aro-dnsmasq-pre.sh
          ExecStart=/usr/sbin/dnsmasq -k
          ExecStopPost=/bin/bash -c '/bin/rm /etc/NetworkManager/conf.d/aro-dns.conf && /usr/bin/nmcli general reload conf && /usr/bin/nmcli general reload dns-rc'
          Restart=always
          StandardOutput=journal+console
          StandardError=journal+console

          [Install]
          WantedBy=multi-user.target
        enabled: true
        name: dnsmasq.service
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`,
	"/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-etc-hosts-gateway-domains.yaml": `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-aro-etc-hosts-gateway-domains
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,MjAzLjAuMTEzLjEJYXBpLnRlc3QtY2x1c3Rlci50ZXN0LmV4YW1wbGUuY29tIGFwaS1pbnQudGVzdC1jbHVzdGVyLnRlc3QuZXhhbXBsZS5jb20KMjAzLjAuMTEzLjIJZ2F0ZXdheS5tb2NrMS5leGFtcGxlLmNvbSBnYXRld2F5Lm1vY2syLmV4YW1wbGUuY29tCg==
        mode: 420
        overwrite: true
        path: /etc/hosts.d/aro.conf
        user:
          name: root
      - contents:
          source: data:text/plain;charset=utf-8;base64,IyEvYmluL2Jhc2gKc2V0IC11byBwaXBlZmFpbAoKdHJhcCAnam9icyAtcCB8IHhhcmdzIGtpbGwgfHwgdHJ1ZTsgd2FpdDsgZXhpdCAwJyBURVJNCgpPUEVOU0hJRlRfTUFSS0VSPSJvcGVuc2hpZnQtYXJvLWV0Y2hvc3RzLXJlc29sdmVyIgpIT1NUU19GSUxFPSIvZXRjL2hvc3RzIgpDT05GSUdfRklMRT0iL2V0Yy9ob3N0cy5kL2Fyby5jb25mIgpURU1QX0ZJTEU9Ii9ldGMvaG9zdHMuZC9hcm8udG1wIgoKIyBNYWtlIGEgdGVtcG9yYXJ5IGZpbGUgd2l0aCB0aGUgb2xkIGhvc3RzIGZpbGUncyBkYXRhLgppZiAhIGNwIC1mICIke0hPU1RTX0ZJTEV9IiAiJHtURU1QX0ZJTEV9IjsgdGhlbgogIGVjaG8gIkZhaWxlZCB0byBwcmVzZXJ2ZSBob3N0cyBmaWxlLiBFeGl0aW5nLiIKICBleGl0IDEKZmkKCmlmICEgc2VkIC0tc2lsZW50ICIvIyAke09QRU5TSElGVF9NQVJLRVJ9L2Q7IHcgJHtURU1QX0ZJTEV9IiAiJHtIT1NUU19GSUxFfSI7IHRoZW4KICAjIE9ubHkgY29udGludWUgcmVidWlsZGluZyB0aGUgaG9zdHMgZW50cmllcyBpZiBpdHMgb3JpZ2luYWwgY29udGVudCBpcyBwcmVzZXJ2ZWQKICBzbGVlcCA2MCAmIHdhaXQKICBjb250aW51ZQpmaQoKd2hpbGUgSUZTPSByZWFkIC1yIGxpbmU7IGRvCiAgICBlY2hvICIke2xpbmV9ICMgJHtPUEVOU0hJRlRfTUFSS0VSfSIgPj4gIiR7VEVNUF9GSUxFfSIKZG9uZSA8ICIke0NPTkZJR19GSUxFfSIKCiMgUmVwbGFjZSAvZXRjL2hvc3RzIHdpdGggb3VyIG1vZGlmaWVkIHZlcnNpb24gaWYgbmVlZGVkCmNtcCAiJHtURU1QX0ZJTEV9IiAiJHtIT1NUU19GSUxFfSIgfHwgY3AgLWYgIiR7VEVNUF9GSUxFfSIgIiR7SE9TVFNfRklMRX0iCiMgVEVNUF9GSUxFIGlzIG5vdCByZW1vdmVkIHRvIGF2b2lkIGZpbGUgY3JlYXRlL2RlbGV0ZSBhbmQgYXR0cmlidXRlcyBjb3B5IGNodXJuCg==
        mode: 484
        overwrite: true
        path: /usr/local/bin/aro-etchosts-resolver.sh
        user:
          name: root
    systemd:
      units:
      - contents: |
          [Unit]
          Description=One shot service that appends static domains to etchosts
          Before=network-online.target

          [Service]
          # ExecStart will copy the hosts defined in /etc/hosts.d/aro.conf to /etc/hosts
          ExecStart=/bin/bash /usr/local/bin/aro-etchosts-resolver.sh

          [Install]
          WantedBy=multi-user.target
        enabled: true
        name: aro-etchosts-resolver.service
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`,
	"/usr/local/bin/aro-dnsmasq-pre.sh": `
#!/bin/bash
set -euo pipefail

# This bash script is a part of the ARO DnsMasq configuration
# It's deployed as part of the 99-aro-dns-* machine config
# See https://github.com/Azure/ARO-RP

# This file can be rerun and the effect is idempotent, output might change if the DHCP configuration changes

NODEIP=$(/sbin/ip --json route get 168.63.129.16 | /bin/jq -r ".[].prefsrc")

if [ "$NODEIP" != "" ]; then
    /bin/cp -Z /etc/resolv.conf /etc/resolv.conf.dnsmasq
    SEARCHDOMAIN=$(awk '/^search/ { print $2; }' /etc/resolv.conf.dnsmasq)
    /bin/chmod 0744 /etc/resolv.conf.dnsmasq

    cat <<EOF | /bin/tee /etc/NetworkManager/conf.d/aro-dns.conf
# Added by dnsmasq.service
[global-dns]
searches=$SEARCHDOMAIN

[global-dns-domain-*]
servers=$NODEIP
EOF

    # network manager may already be running at this point.
    # reload to update /etc/resolv.conf with this configuration
    /usr/bin/nmcli general reload conf
    /usr/bin/nmcli general reload dns-rc
fi
`,
	"/usr/local/bin/aro-etchosts-resolver.sh": `#!/bin/bash
set -uo pipefail

trap 'jobs -p | xargs kill || true; wait; exit 0' TERM

OPENSHIFT_MARKER="openshift-aro-etchosts-resolver"
HOSTS_FILE="/etc/hosts"
CONFIG_FILE="/etc/hosts.d/aro.conf"
TEMP_FILE="/etc/hosts.d/aro.tmp"

# Make a temporary file with the old hosts file's data.
if ! cp -f "${HOSTS_FILE}" "${TEMP_FILE}"; then
  echo "Failed to preserve hosts file. Exiting."
  exit 1
fi

if ! sed --silent "/# ${OPENSHIFT_MARKER}/d; w ${TEMP_FILE}" "${HOSTS_FILE}"; then
  # Only continue rebuilding the hosts entries if its original content is preserved
  sleep 60 & wait
  continue
fi

while IFS= read -r line; do
    echo "${line} # ${OPENSHIFT_MARKER}" >> "${TEMP_FILE}"
done < "${CONFIG_FILE}"

# Replace /etc/hosts with our modified version if needed
cmp "${TEMP_FILE}" "${HOSTS_FILE}" || cp -f "${TEMP_FILE}" "${HOSTS_FILE}"
# TEMP_FILE is not removed to avoid file create/delete and attributes copy churn
`,
	"/opt/openshift/manifests/cluster-dns-02-config.yml": `apiVersion: config.openshift.io/v1
kind: DNS
metadata:
  creationTimestamp: null
  name: cluster
spec:
  baseDomain: test-cluster.test.example.com
  platform:
    aws: null
    type: ""
status: {}
`,
	"/opt/openshift/openshift/99_openshift-cluster-api_master-user-data-secret.yaml": `apiVersion: v1
kind: Secret
metadata:
  name: master-user-data
  namespace: openshift-machine-api
type: Opaque
data:
  disableTemplating: "dHJ1ZQo="
  userData: test
`,
	"/opt/openshift/openshift/99_openshift-cluster-api_worker-user-data-secret.yaml": `apiVersion: v1
kind: Secret
metadata:
  name: worker-user-data
  namespace: openshift-machine-api
type: Opaque
data:
  disableTemplating: "dHJ1ZQo="
  userData: test
`,
}
