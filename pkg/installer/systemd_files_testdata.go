package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

var expectedIgnitionServiceContents = map[string]string{
	"aro-etchosts-resolver.service": `[Unit]
Description=One shot service that appends static domains to etchosts
Before=node-image-pull.service

[Service]
Type=oneshot
# ExecStart will copy the hosts defined in /etc/hosts.d/aro.conf to /etc/hosts
ExecStart=/bin/bash /usr/local/bin/aro-etchosts-resolver.sh

[Install]
WantedBy=network-online.target
`,
	"dnsmasq.service": `
[Unit]
Description=DNS caching server.
After=network-online.target
Wants=network-online.target
Before=bootkube.service
Before=node-image-pull.service

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
WantedBy=network.target
`,
}
